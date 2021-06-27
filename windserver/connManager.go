package windserver

import (
	"github.com/panjf2000/gnet"
	"log"
	"time"
)


type ConnManager struct {
	*gnet.EventServer
	srv 		*windServer
}


func NewConnManager(inst *windServer)  *ConnManager {
	return &ConnManager{srv: inst}
}

func (cm *ConnManager) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	// Echo asynchronously.
	data := append([]byte{}, frame...)
	go func() {
		time.Sleep(time.Second)
		c.AsyncWrite(data)
	}()
	return
}

func (cm *ConnManager) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("Echo server is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

func (cm *ConnManager) OnOpened(c gnet.Conn) (action gnet.Action) {
	return
}

func (cm *ConnManager) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	return
}



