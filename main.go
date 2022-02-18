package main

import (
	"flag"
	"fmt"
	"github.com/F3eev/gobfroce/lib/nmapxml"
	"github.com/F3eev/gobfroce/scan"
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
		if hostInfo.Status.State == "down" {
			fmt.Printf("%s is not online\n", hostInfo.Address.Addr)
			break
		}
		for _, port := range *hostInfo.Ports.Port {
			result = append(result, nmapResult{ip, port.PortID, port.Service.Name})
		}
	}
	return result, nil
}

func main() {

	argThreads := flag.Int("threads", 400, "thread num default 400")
	argNmapFile := flag.String("nFile", "", "nmap xml file")
	argTimeOut := flag.Int("timeout", 5, "timeout")
	argNmapDir := flag.String("nDir", "", "nmap xml file")
	argOnlyCustomDict := flag.Bool("CustomDict", false, "only use *_custom dict (default false)")
	argSelectService := flag.String("service", "all", "choose service to scan")
	argLog := flag.String("log", "log.txt", "log file")
	argLevel := flag.Int("level", 2, "log file")

	flag.Parse()

	var nmapResultXmlList []nmapResult
	var nmapOutFileList []string
	var targets []scan.Target
	if *argNmapDir != "" {
		fileInfoList, _ := ioutil.ReadDir(*argNmapDir)
		for i := range fileInfoList {
			nmapOutFileList = append(nmapOutFileList, *argNmapDir+"/"+fileInfoList[i].Name())
		}
	}
	if *argNmapFile != "" {
		nmapOutFileList = append(nmapOutFileList, *argNmapFile)
	}
	for _, f := range nmapOutFileList {
		if temp, err := loadNmap(f); err == nil {
			nmapResultXmlList = append(nmapResultXmlList, temp...)
		} else {
			fmt.Printf("%s %s \n", f, err.Error())
		}
	}

	for _, t := range nmapResultXmlList {
		targets = append(targets, scan.Target{IP: t.IP, Port: t.Port, Service: t.Service})
	}
	if len(targets) == 0 {
		flag.Usage()
		return
	}

	v := scan.Init(*argThreads, *argSelectService, *argOnlyCustomDict, *argLog, *argTimeOut, *argLevel)
	//targets = append(targets, scan.Target{IP: "127.0.0.1", Port: "27017", Service: "mongodb"})
	//
	//targets = append(targets, scan.Target{IP: "127.0.0.1", Port: "277", Service: "mongodb"})
	////
	//targets = append(targets, scan.Target{IP: "127.0.0.1", Port: "27017", Service: "x"})
	v.BruteForce(targets)

}
