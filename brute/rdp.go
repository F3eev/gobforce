package brute

import (

	"github.com/icodeface/grdp"
	"github.com/icodeface/grdp/glog"
)

func (t *Target)RdpLogin() bool {


	client := grdp.NewClient(t.IP+":"+t.Port, glog.NONE)
	err := client.Login(t.Username, t.Password)
	if err != nil {
		return false
	} else {
		return true
	}
}