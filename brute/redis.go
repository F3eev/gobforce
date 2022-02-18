package brute

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

func (t *Target) RedisLogin() bool {

	debugMessage := fmt.Sprintf("%s %s %s:%s", RunFileName(), RunFuncName(), t.IP, t.Port)

	c, err := redis.Dial("tcp", t.IP+":"+t.Port)
	if err != nil {
		t.Logrus.Debug(fmt.Sprintf("%s %v", debugMessage, err))

		return false
	}
	err = c.Send("auth", t.Password)
	if err != nil {
		t.Logrus.Debug(fmt.Sprintf("%s %v", debugMessage, err))

		return false
	}
	return true
}
