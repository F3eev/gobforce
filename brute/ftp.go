package brute

import (
	"fmt"
	"github.com/prasad83/goftp"
	"time"
)

func (t *Target) FTPLogin() bool {

	configWithoutTLS := goftp.Config{
		User:               t.Username,
		Password:           t.Password,
		ConnectionsPerHost: 5,
		Timeout:            10 * time.Second,
	}
	debugMessage := fmt.Sprintf("%s %s %s:%s", RunFileName(), RunFuncName(), t.IP, t.Port)
	client, err := goftp.DialConfig(configWithoutTLS, fmt.Sprintf("%s:%s", t.IP, t.Port))
	if err != nil {
		t.Logrus.Debug(fmt.Sprintf("%s %v", debugMessage, err))

		return false
	}
	if _, err := client.ReadDir("/"); err != nil {
		client.Close()
		t.Logrus.Debug(fmt.Sprintf("%s %v", debugMessage, err))
		return false
	}
	defer client.Close()
	return true
}
