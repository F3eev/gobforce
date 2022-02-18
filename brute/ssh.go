package brute

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

const chacha20Poly1305ID = "chacha20-poly1305@openssh.com"
const (
	gcmCipherID    = "aes128-gcm@openssh.com"
	aes128cbcID    = "aes128-cbc"
	tripledescbcID = "3des-cbc"
)

var supportedCiphers = []string{
	"aes128-ctr", "aes192-ctr", "aes256-ctr",
	chacha20Poly1305ID,
	"arcfour256", "arcfour128", "arcfour",
	aes128cbcID,
	tripledescbcID,
	gcmCipherID,
}

func (t *Target) SSHLogin() bool {
	debugMessage := fmt.Sprintf("%s %s %s:%s", RunFileName(), RunFuncName(), t.IP, t.Port)

	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(t.Password))

	clientConfig = &ssh.ClientConfig{
		User:    t.Username,
		Auth:    auth,
		Timeout: 4 * time.Second,
		//2019.6.18  golang默认配置加密方式不包括aes128-cbc  连接交换机需要使用aes128-cbc
		Config: ssh.Config{
			Ciphers: supportedCiphers,
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%s", t.IP, t.Port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		t.Logrus.Debug(fmt.Sprintf("%s %v", debugMessage, err))
		return false
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		t.Logrus.Debug(fmt.Sprintf("%s %v", debugMessage, err))
		return false
	}
	session.Close()
	return true
}
