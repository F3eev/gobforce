package brute

import (
	"fmt"
	"github.com/prasad83/goftp"
	"log"

	"time"
)

func (t *Target)FTPLogin() bool {

	configWithoutTLS := goftp.Config{
		User:               t.Username,
		Password:           t.Password,
		ConnectionsPerHost: 5,
		Timeout:            10 * time.Second,
	}

	client, err := goftp.DialConfig(configWithoutTLS, fmt.Sprintf("%s:%s", t.IP, t.Port))
	if err !=nil{
		log.Println(err.Error())
		return false
	}
	if _, err := client.ReadDir("/"); err != nil {
		client.Close()
		log.Println(err.Error())
		return false
	}
	defer client.Close()
	return true
}