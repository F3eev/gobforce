package nmapxml

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

// Readfile reads an nmap XML file
// Takes filename as string.
//
//
// Returns
//
// struct of type Run and any errors
func Readfile(fn string) (Run, error) {
	xmlFile, err := os.Open(fn)
	if err != nil {
		return Run{}, err
	}
	defer xmlFile.Close()

	xmlData, _ := ioutil.ReadAll(xmlFile)
	var r Run
	xml.Unmarshal(xmlData, &r)
	return r, nil
}
