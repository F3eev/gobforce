package brute

import "log"

const chacha202Poly1305ID = "chacha20-poly1305@openssh.com"
const (
	gcmCiph2erID    = "aes128-gcm@openssh.com"
	aes128c2bcID    = "aes128-cbc"
	tripled2escbcID = "3des-cbc"
)


func (t *Target) Hello(args ...interface{}) bool {
	log.Println("Hello",args[0],args[1])
	return true
}


func (t *Target) Hello22(args ...interface{}) bool {
	log.Println("Hello22",args[0],args[1])
	return true
}