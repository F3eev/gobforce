package brute

import (
	"context"
	"fmt"
	"github.com/kward/go-vnc"
	"net"
)

func (t *Target) VNCLoginNoUser() bool {

	debugMessage := fmt.Sprintf("%s %s %s:%s", RunFileName(), RunFuncName(), t.IP, t.Port)
	nc, err := net.Dial("tcp", t.IP+":"+t.Port)
	if err != nil {
		t.Logrus.Debug(fmt.Sprintf("%s %v", debugMessage, err))
		return false
	}
	vcc := vnc.NewClientConfig(t.Password)
	_, err = vnc.Connect(context.Background(), nc, vcc)
	if err != nil {
		t.Logrus.Debug(fmt.Sprintf("%s %v", debugMessage, err))
		return false
	}
	return true

}
