package brute

import "log"

func (t *Target) Hello(args ...interface{}) bool {
	log.Println("Hello", args[0], args[1])
	return true
}

func (t *Target) Hello22(args ...interface{}) bool {
	log.Println("Hello22", args[0], args[1])
	return true
}
