package brute

import (
	"fmt"
	"github.com/tomatome/grdp/core"
	"github.com/tomatome/grdp/glog"
	"github.com/tomatome/grdp/protocol/nla"
	"github.com/tomatome/grdp/protocol/pdu"
	"github.com/tomatome/grdp/protocol/rfb"
	"github.com/tomatome/grdp/protocol/sec"
	"github.com/tomatome/grdp/protocol/t125"
	"github.com/tomatome/grdp/protocol/tpkt"
	"github.com/tomatome/grdp/protocol/x224"
	"log"
	"net"
	"os"
	"sync"
)

func (t *Target) RdpLogin() bool {

	target := fmt.Sprintf("%s:%s", t.IP, t.Port)
	g := NewClient(target, glog.NONE)
	debugMessage := fmt.Sprintf("%s %s %s:%s", RunFileName(), RunFuncName(), t.IP, t.Port)
	conn, err := net.Dial("tcp", g.Host)
	if err != nil {
		t.Logrus.Debug(fmt.Sprintf("%s %v", debugMessage, err))
		return false
	}
	g.tpkt = tpkt.New(core.NewSocketLayer(conn), nla.NewNTLMv2("domain", t.Username, t.Password))
	g.x224 = x224.New(g.tpkt)
	g.mcs = t125.NewMCSClient(g.x224)
	g.sec = sec.NewClient(g.mcs)
	g.pdu = pdu.NewClient(g.sec)

	g.sec.SetUser(t.Username)
	g.sec.SetPwd(t.Password)
	g.sec.SetDomain("domain")
	//g.sec.SetClientAutoReconnect()

	g.tpkt.SetFastPathListener(g.sec)
	g.sec.SetFastPathListener(g.pdu)
	g.pdu.SetFastPathSender(g.tpkt)

	//g.x224.SetRequestedProtocol(x224.PROTOCOL_SSL)
	//g.x224.SetRequestedProtocol(x224.PROTOCOL_RDP)

	err = g.x224.Connect()
	if err != nil {
		t.Logrus.Debug(fmt.Sprintf("%s %v", debugMessage, err))
		return false
	}
	wg := &sync.WaitGroup{}
	breakFlag := false
	wg.Add(1)
	g.pdu.On("error", func(e error) {
		err = e
		t.Logrus.Debug(fmt.Sprintf("%s %v", debugMessage, err))
		g.pdu.Emit("done")
	})
	g.pdu.On("close", func() {

		g.pdu.Emit("done")
	})
	g.pdu.On("success", func() {
		err = nil
		g.pdu.Emit("done")
	})
	g.pdu.On("ready", func() {

		g.pdu.Emit("done")
	})
	g.pdu.On("update", func(rectangles []pdu.BitmapData) {

	})
	g.pdu.On("done", func() {
		if breakFlag == false {
			breakFlag = true
			wg.Done()
		}
	})

	wg.Wait()

	//var err
	if err != nil {
		t.Logrus.Debug(fmt.Sprintf("%s %v", debugMessage, err))
		return false
	}
	//return true, err
	return true
}

type Client struct {
	Host string // ip:port
	tpkt *tpkt.TPKT
	x224 *x224.X224
	mcs  *t125.MCSClient
	sec  *sec.Client
	pdu  *pdu.Client
	vnc  *rfb.RFB
}

func NewClient(host string, logLevel glog.LEVEL) *Client {
	glog.SetLevel(logLevel)
	logger := log.New(os.Stdout, "", 0)
	glog.SetLogger(logger)
	return &Client{
		Host: host,
	}
}
