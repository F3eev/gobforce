package scan

import (
	"github.com/F3eev/gobfroce/brute"
	"github.com/F3eev/gobfroce/lib"
	"github.com/panjf2000/ants"
	"io"
	"log"
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
	Logger              *log.Logger
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
}

var expList = map[string][]string{
	"ssh":        {"SSHLogin"},
	"ftp":        {"FTPLogin"},
	"mongodb":    {"MongoLogin"},
	"mysql":      {"MysqlLogin"},
	"postgresql": {"PostgresLogin"},
	//"ms-wbt-server": {"RdpLogin"},
	"redis": {"RedisLogin"},
	"vnc":   {"VNCLoginNoUser"},
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

func Init(threadNum int, services string, onlyCustomDict bool, logfile string, timeout int) *VScan {
	logFile, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		log.Printf(":Logfile %s %s\n", logfile, err.Error())
	}
	//defer logFile.Close()
	fileAndPrint := io.MultiWriter(logFile, os.Stdout)
	logger := log.New(fileAndPrint, "[BruteForce]", log.LstdFlags)
	serviceDictFunction := make(map[string]ServiceDictFunction)
	scanService := filterService(services)
	logger.Printf(":Staring scan thread %d can service:%v ", threadNum, scanService)

	for _, service := range scanService {
		for _, fun := range expList[service] {
			dict := loadExpDict(service, onlyCustomDict)
			logger.Printf(":Load %s dict %d", service, len(dict))
			factor := ServiceDictFunction{service, dict, fun}
			serviceDictFunction[service] = factor
		}
	}

	return &VScan{serviceDictFunction, scanService, threadNum, logger, timeout}
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

func checkOpen(target Target) bool {

	connection, err := net.DialTimeout("tcp", target.IP+":"+target.Port, 5*time.Second)
	if err != nil {
		return false
	}
	defer connection.Close()
	return true
}

func (v *VScan) BruteForce(targets []Target) {

	t1 := time.Now()

	v.Logger.Printf(":Target:%d", len(targets))
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(v.threadNum, func(i interface{}) {
		scanPar := i.(ScanPar)
		bruteExp := brute.Target{IP: scanPar.Target.IP, Port: scanPar.Target.Port, Username: scanPar.Username, Password: scanPar.Password}
		result, err := bruteExp.CallFunc(scanPar.Function)
		if err == nil {
			status := result[0].Bool()
			if status == true {
				//	success [1;40;32m0x1B0x1B
				v.Logger.Printf(":%c[1;40;32m%s %s %s:%s  (%s:%s) Success%c[0m", 0x1B, scanPar.Target.Service, scanPar.Function, scanPar.Target.IP, scanPar.Target.Port, scanPar.Username, scanPar.Password, 0x1B)

			} else {
				v.Logger.Printf(":%s %s %s:%s (%s:%s) Fail", scanPar.Target.Service, scanPar.Function, scanPar.Target.IP, scanPar.Target.Port, scanPar.Username, scanPar.Password)
			}
		}
		wg.Done()
	})
	defer p.Release()

	for _, target := range targets {
		if checkOpen(target) {
			_, ok := v.serviceDictFunction[target.Service]
			if ok {
				for _, userPass := range v.serviceDictFunction[target.Service].DictUserPass {
					scanPar := ScanPar{target, v.serviceDictFunction[target.Service].Function, userPass.Username, userPass.Password}
					wg.Add(1)
					_ = p.Invoke(scanPar)
				}
			} else {
				v.Logger.Printf(":%s %s:%s is unable to hit", target.Service, target.IP, target.Port)
			}
		} else {
			v.Logger.Printf(":%s %s:%s is not open", target.Service, target.IP, target.Port)
		}
	}
	wg.Wait()
	elapsed := time.Since(t1)
	v.Logger.Printf("Finish scan %s", elapsed)
}
