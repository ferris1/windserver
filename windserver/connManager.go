package windserver

import (
	"github.com/panjf2000/gnet"
)

type ClientConn struct {
	conn 			*gnet.Conn
	peerId			int
}

func NewClientConn(cn *gnet.Conn, id int)  *ClientConn{
	return &ClientConn{conn: cn, peerId: id}
}

// 代理一层网络转发  解耦服务器连接管理与网络管理的逻辑，方便之后替换网络层
type ConnManager struct {
	srv 				*windServer
	peerToClient		map[int]*ClientConn
	totalConnNumber     int
}

func NewConnManager(inst *windServer)  *ConnManager {
	mgr:= &ConnManager{srv: inst}
	mgr.peerToClient = make(map[int]*ClientConn)
	mgr.totalConnNumber = 0
	return mgr
}

func (cm *ConnManager) OnConnect(peerId int, conn *gnet.Conn) {
	cm.peerToClient[peerId] = NewClientConn(conn, peerId)
	cm.totalConnNumber ++
}

func (cm *ConnManager) OnDisConnect(peerId int)  {
	delete(cm.peerToClient, peerId)
}

func (cm *ConnManager) OnData(peerId int, frame []byte)  {

}


