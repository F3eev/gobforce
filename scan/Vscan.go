package scan

import (
	"bytes"
	"fmt"
	"github.com/F3eev/gobfroce/brute"
	"github.com/F3eev/gobfroce/lib"
	"github.com/cheggaaa/pb/v3"
	"github.com/panjf2000/ants"
	logrus "github.com/sirupsen/logrus"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type VScan struct {
	serviceDictFunction map[string]ServiceDictFunction
	serviceAction       []string
	threadNum           int
	Logrus              *logrus.Logger
	timeout             int
}

type Target struct {
	IP      string
	Port    string
	Service string
}

type DictUserPass struct {
	Username string
	Password string
}
type ServiceDictFunction struct {
	Service      string
	DictUserPass []DictUserPass
	Function     string
}
type ScanPar struct {
	Target   Target
	Function string
	Username string
	Password string
	Logrus   *logrus.Logger
}

var expList = map[string][]string{
	"ssh":           {"SSHLogin"},
	"ftp":           {"FTPLogin"},
	"mongodb":       {"MongoLogin"},
	"mysql":         {"MysqlLogin"},
	"postgresql":    {"PostgresLogin"},
	"ms-wbt-server": {"RdpLogin"},
	"redis":         {"RedisLogin"},
	"vnc":           {"VNCLoginNoUser"},
}

func getDefaultService() (services []string) {
	for s, _ := range expList {
		services = append(services, s)
	}

	return services
}

func filterService(ArgServices string) []string {
	var inputServiceList []string
	var scanServiceList []string

	if strings.Contains(ArgServices, ",") {
		inputServiceList = strings.Split(ArgServices, ",")
	} else {
		if ArgServices == "all" {
			inputServiceList = getDefaultService()
			scanServiceList = inputServiceList
			return inputServiceList
		} else {
			inputServiceList = append(inputServiceList, ArgServices)
		}
	}
	for _, service := range inputServiceList {
		_, ok := expList[service]
		if ok {
			scanServiceList = append(scanServiceList, service)
		}
	}
	return scanServiceList
}

type MyFormatter struct {
	Prefix string
	Suffix string
}

func (mf *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := time.Now().Local().Format("2006/1/02 15:04")

	b.WriteString(fmt.Sprintf("%s%s:  - %s\n", mf.Prefix, timestamp, entry.Message))
	return b.Bytes(), nil
}

func Init(threadNum int, services string, onlyCustomDict bool, logfile string, timeout int, level int) *VScan {

	var levelFlag logrus.Level
	if level == 2 {
		levelFlag = logrus.DebugLevel
	} else if level == 1 {
		levelFlag = logrus.InfoLevel
	}

	log := logrus.New()
	log.SetLevel(levelFlag)
	formatter := &MyFormatter{
		Prefix: "[gobfroce]",
	}
	log.SetFormatter(formatter)
	logFile, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		log.Error("Logfile %s %s\n", logfile, err.Error())
	}
	log.SetOutput(io.MultiWriter(logFile, os.Stdout))
	serviceDictFunction := make(map[string]ServiceDictFunction)
	scanService := filterService(services)
	log.Printf("Staring scan thread %d can service:%v ", threadNum, scanService)

	for _, service := range scanService {
		for _, fun := range expList[service] {
			dict := loadExpDict(service, onlyCustomDict)
			log.Printf("Load %s dict %d", service, len(dict))
			factor := ServiceDictFunction{service, dict, fun}
			serviceDictFunction[service] = factor
		}
	}

	return &VScan{serviceDictFunction, scanService, threadNum, log, timeout}
}

func loadExpDict(name string, onlyCustomDict bool) []DictUserPass {
	var dictUserPass []DictUserPass
	commonUser := lib.FileRead("dict/common_user.txt")
	commonPass := lib.FileRead("dict/common_pass.txt")
	commonCustom := lib.FileRead("dict/common_custom.txt")
	serviceCustom := lib.FileRead("dict/" + name + "_custom.txt")
	serverUser := lib.FileRead("dict/" + name + "_user.txt")
	serverPass := lib.FileRead("dict/" + name + "_password.txt")

	userAll := lib.RemoveDuplicateElement(append(commonUser, serverUser...))
	passAll := lib.RemoveDuplicateElement(append(commonPass, serverPass...))
	commonAll := lib.RemoveDuplicateElement(append(commonCustom, serviceCustom...))

	for _, line := range commonAll {
		arr := strings.Split(line, ":")
		dictUserPass = append(dictUserPass, DictUserPass{arr[0], arr[1]})
	}
	// if custom, only return custom dict
	if onlyCustomDict {
		return dictUserPass
	}
	for _, user := range userAll {
		for _, pass := range passAll {
			dictUserPass = append(dictUserPass, DictUserPass{user, pass})
		}
	}
	for _, line := range commonCustom {
		arr := strings.Split(line, ":")
		dictUserPass = append(dictUserPass, DictUserPass{arr[0], arr[1]})
	}

	return dictUserPass
}

func (v *VScan) checkOpen(target Target) bool {

	connection, err := net.DialTimeout("tcp", target.IP+":"+target.Port, 5*time.Second)
	if err != nil {
		return false
	}
	defer connection.Close()
	return true
}

func (v *VScan) BruteForce(targets []Target) {

	t1 := time.Now()
	bar := pb.New(len(targets))
	bar.Start()
	v.Logrus.Printf("Target:%d", len(targets))
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(v.threadNum, func(i interface{}) {
		scanPar := i.(ScanPar)
		bruteExp := brute.Target{IP: scanPar.Target.IP, Port: scanPar.Target.Port, Username: scanPar.Username, Password: scanPar.Password, Logrus: scanPar.Logrus}
		result, err := bruteExp.CallFunc(scanPar.Function)
		if err == nil {
			status := result[0].Bool()
			if status == true {
				v.Logrus.Info(fmt.Sprintf("%c[1;40;32m%s %s %s:%s (%s:%s) successful %c[0m", 0x1B, scanPar.Target.Service, scanPar.Function, scanPar.Target.IP, scanPar.Target.Port, scanPar.Username, scanPar.Password, 0x1B))
			} else {
				message := fmt.Sprintf("%s %s %s:%s (%s:%s) Fail", scanPar.Target.Service, scanPar.Function, scanPar.Target.IP, scanPar.Target.Port, scanPar.Username, scanPar.Password)
				v.Logrus.Info(message)
			}
		}

		wg.Done()
	})
	defer p.Release()

	for _, target := range targets {
		bar.Increment()
		if v.checkOpen(target) {
			_, ok := v.serviceDictFunction[target.Service]
			if ok {
				v.Logrus.Printf("Start brute %s:%s %s ", target.IP, target.Port, target.Service)

				for _, userPass := range v.serviceDictFunction[target.Service].DictUserPass {
					scanPar := ScanPar{target, v.serviceDictFunction[target.Service].Function, userPass.Username, userPass.Password, v.Logrus}
					wg.Add(1)
					_ = p.Invoke(scanPar)
				}
			} else {
				v.Logrus.Printf("%s %s:%s is unable to hit", target.Service, target.IP, target.Port)
			}
		} else {
			v.Logrus.Printf("%s %s:%s is not open", target.Service, target.IP, target.Port)
		}
	}
	wg.Wait()
	bar.Finish()
	elapsed := time.Since(t1)
	v.Logrus.Printf("Finish scan %s", elapsed)
}
