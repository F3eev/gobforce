package brute

import "reflect"

type Target struct {
	IP string
	Port string
	Username string
	Password string
}

type FuncCollection map[string]reflect.Value

func (t Target)CallFunc(tableName string, args ... interface{}) (result []reflect.Value, err error) {
	//var exp Exp
	FuncMap := make(FuncCollection, 0)
	rf := reflect.ValueOf(&t)
	rft := rf.Type()
	funcNum := rf.NumMethod()
	for i := 0; i < funcNum; i ++ {
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