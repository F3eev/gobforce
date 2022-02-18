package lib

import (
	"bufio"
	"os"
	"strings"
)

func FileRead(fileName string) []string {
	var fileContent []string

	if file, err := os.Open(fileName); err != nil {
		panic(err)
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fileContent = append(fileContent, scanner.Text())
		}
	}
	return fileContent
}

func RemoveDuplicateElement(source []string) []string {
	result := make([]string, 0, len(source))
	temp := map[string]struct{}{}
	for _, item := range source {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, strings.Replace(item, "\n", "", -1))
		}
	}
	return result
}
