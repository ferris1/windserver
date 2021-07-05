package windserver

import (
	"context"
	"github.com/panjf2000/gnet"
)

type NetServer struct {
	unusedPeerId		int
	connToPeerId        map[*gnet.Conn]int
	connMgr				*ConnManager
}

func NewNetServer(connMgr *ConnManager) *NetServer{
	ns := &NetServer{ unusedPeerId: 0}
	ns.connToPeerId = make(map[*gnet.Conn]int)
	ns.connMgr = connMgr
	return ns
}

func (ns *NetServer) StartService(ctx context.Context, ip string, port int, isTcp bool)  error {
	return nil
}



func (ns *NetServer) AllocNewPeerId() int {
	ns.unusedPeerId ++
	return ns.unusedPeerId
}