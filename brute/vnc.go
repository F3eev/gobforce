package brute

import (
	"context"
	"github.com/kward/go-vnc"
	"net"
)

func (t *Target) VNCLoginNoUser() bool {

	// Establish TCP connection to VNC server.
	nc, err := net.Dial("tcp", t.IP+":"+t.Port)
	if err != nil {
		// Negotiate connection with the server.
		//log.Fatalf("Error connecting to VNC host %s %s. %v",t.IP,t.Port, err)
		return false

	}
	vcc := vnc.NewClientConfig(t.Password)
	_, err = vnc.Connect(context.Background(), nc, vcc)
	if err != nil {
		return false
	}

	return true

}