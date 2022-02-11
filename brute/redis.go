package brute

import (
	"github.com/garyburd/redigo/redis"

)

func (t *Target) RedisLogin() bool {


	c, err := redis.Dial("tcp", t.IP+":"+t.Port)
	if err != nil {
		return false
	}
	err = c.Send("auth", t.Password)
	if err != nil {
		return false
	}
	return true
}

