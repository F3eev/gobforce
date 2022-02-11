package main

import (
	"Goscanpro/lib/nmapxml"
	"Goscanpro/scan"
	"flag"
	"fmt"
	"io/ioutil"
)

type nmapResult struct {
	IP      string
	Port    string
	Service string
}

func loadNmap(file string) ([]nmapResult, error) {

	var result []nmapResult
	scanData, _ := nmapxml.Readfile(file)
	for _, hostInfo := range scanData.Host {
		ip := hostInfo.Address.Addr
		if hostInfo.Status.State =="down"{
			fmt.Printf("%s is not online\n",hostInfo.Address.Addr)
			break
		}
		for _, port := range *hostInfo.Ports.Port {
			result = append(result, nmapResult{ip, port.PortID, port.Service.Name})
		}
	}
	return result, nil
}

func main() {


	argThreads := flag.Int("threads", 500, "thread num default 400")
	argNmapFile := flag.String("nFile", "", "nmap xml file")
	argTimeOut := flag.Int("-to", 5, "timeout")
	argNmapDir := flag.String("nDir", "nmapOutXml", "nmap xml file")
	argDict := flag.String("dict", "", "custom only for *_custom dict")
	argSelectService := flag.String("service", "all", "selecting service to scan")
	argLog := flag.String("log", "log.txt", "log file")

	flag.Parse()
	var nmapResultXmlList []nmapResult
	var nmapOutFileList []string
	var targets []scan.Target
	if *argNmapDir != "" {
		fileInfoList, _ := ioutil.ReadDir(*argNmapDir)
		for i := range fileInfoList {
			nmapOutFileList = append(nmapOutFileList, *argNmapDir + "/" + fileInfoList[i].Name())
		}
	}
	if *argNmapFile != "" {
		nmapOutFileList = append(nmapOutFileList,*argNmapFile )
	}
	for _, f := range nmapOutFileList {
		if temp, err:=loadNmap(f);err ==nil{
			nmapResultXmlList = append(nmapResultXmlList, temp...)
		} else{
			fmt.Printf("%s %s \n",f,err.Error())
		}
	}

	for _, t := range nmapResultXmlList {
		if t.Service =="ftp"{
			fmt.Println(t.IP, t.Port)
		}
		targets = append(targets, scan.Target{IP: t.IP, Port: t.Port, Service: t.Service})
	}

	v := scan.Init(*argThreads, *argSelectService, *argDict, *argLog,*argTimeOut)
	//targets = append(targets, scan.Target{IP: "47.93.63.243", Port: "22", Service: "ssh"})
	//targets = append(targets, scan.Target{IP: "101.201.123.178", Port: "60021", Service: "ftp"})
	//targets = append(targets, scan.Target{IP: "127.0.0.1", Port: "4567", Service: "ftp"})

	v.BruteForce(targets)

}
