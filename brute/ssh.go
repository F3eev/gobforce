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

//func (t *Target) SSHLogin(args ...interface{}) bool {
//	log.Println(t.IP)
//	//ip:=lib.InterToString(value[0])
//	//port:=lib.InterToString(value[1])
//	//username :=lib.InterToString(value[2])
//	//password :=lib.InterToString(value[3])
//	////println(ip,port,username,password)
//	return true
//}
func (t *Target) SSHLogin() bool {
	//log.Fatal("sshLogin")

	//ip := t.IP
	//port := t.Port
	//username := t.Username
	//password := t.Password
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
		//log.Print("用户名：", username, "    密码: ", password, "      ", err.Error())
		return false
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		//log.Print("用户名：", username, "    密码: ", password, "      ", err.Error())
		return false
	}
	session.Close()
	//log.Print("用户名：", username, "    密码: ", password, "     sssss ")

	//define.Output(value)
	return true
}
