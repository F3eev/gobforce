package brute

import (
	"gopkg.in/mgo.v2"
)

func (t *Target)MongoLogin() bool {
	session, err := mgo.Dial("mongodb://"+t.Username+":"+t.Password+"@"+t.IP+":"+t.Port+"/admin")


	if err != nil {
		return false
	} else {
		defer session.Close()
		return true
	}
}
