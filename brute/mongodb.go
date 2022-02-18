package brute

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

func (t *Target) MongoLogin() bool {
	session, err := mgo.Dial("mongodb://" + t.Username + ":" + t.Password + "@" + t.IP + ":" + t.Port + "/admin")
	debugMessage := fmt.Sprintf("%s %s %s:%s", RunFileName(), RunFuncName(), t.IP, t.Port)

	if err != nil {
		t.Logrus.Debug(fmt.Sprintf("%s %v", debugMessage, err))
		return false
	} else {
		defer session.Close()
		return true
	}
}
