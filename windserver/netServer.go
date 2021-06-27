package windserver

import (
	"context"
	"fmt"
	"github.com/panjf2000/gnet"
	"log"
)

type NetServer struct {
	*gnet.EventServer
	unusedPeerId		int
	connToPeerId        map[*gnet.Conn]int
	connMgr				*ConnManager
}

func NewNetServer(connMgr *ConnManager) *NetServer{
	ns := &NetServer{EventServer:&gnet.EventServer{}, unusedPeerId: 0}
	ns.connToPeerId = make(map[*gnet.Conn]int)
	ns.connMgr = connMgr
	return ns
}

func (ns *NetServer) StartService(ctx context.Context, ip string, port int, isTcp bool)  error {
	netPrefix := ""
	if isTcp {	netPrefix = "tcp"	} else {	netPrefix = "udp"	}
	err := gnet.Serve(ns.EventServer, fmt.Sprintf("%s://:%d", netPrefix, port), gnet.WithMulticore(NETUSEMUlTICORE),
		gnet.WithReusePort(NETREUSEPORT))
	return err
}

func (ns *NetServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	peerId, ok := ns.connToPeerId[&c]
	if ok {
		ns.connMgr.OnData(peerId, frame)
	} else {
		println("error:not found conn:", c.RemoteAddr().String())
	}

	return
}

func (ns *NetServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("Echo server is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

func (ns *NetServer) OnOpened(c gnet.Conn) (action gnet.Action) {
	peerId := ns.AllocNewPeerId()
	ns.connToPeerId[&c] = peerId
	ns.connMgr.OnConnect(peerId, &c)
	return
}

func (ns *NetServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	peerId, ok := ns.connToPeerId[&c]
	if ok {
		ns.connMgr.OnDisConnect(peerId)
	} else {
		println("error:not found conn:", c.RemoteAddr().String())
	}
	delete(ns.connToPeerId, &c)
	return
}

func (ns *NetServer) AllocNewPeerId() int {
	ns.unusedPeerId ++
	return ns.unusedPeerId
}