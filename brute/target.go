package brute

import (
	"github.com/sirupsen/logrus"
	"reflect"
	"runtime"
	"strings"
)

type Target struct {
	IP       string
	Port     string
	Username string
	Password string
	Logrus   *logrus.Logger
}

type FuncCollection map[string]reflect.Value

func (t Target) CallFunc(tableName string, args ...interface{}) (result []reflect.Value, err error) {
	//var exp Exp
	FuncMap := make(FuncCollection, 0)
	rf := reflect.ValueOf(&t)
	rft := rf.Type()
	funcNum := rf.NumMethod()
	for i := 0; i < funcNum; i++ {
		mName := rft.Method(i).Name
		FuncMap[mName] = rf.Method(i)
	}

	parameter := make([]reflect.Value, len(args))
	for k, arg := range args {
		parameter[k] = reflect.ValueOf(arg)
	}
	result = FuncMap[tableName].Call(parameter)
	return
}

func RunFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])

	strArr := strings.Split(f.Name(), ".")
	return strArr[len(strArr)-1]
}
func RunFileName() string {
	_, filePath, _, _ := runtime.Caller(1)
	strArr := strings.Split(filePath, "/")
	file := strArr[len(strArr)-1]
	strArr = strings.Split(file, ".")
	return strArr[0]
}
