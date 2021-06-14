package engine

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"strings"
)

type ServerGroupManager struct {
	etcdClient     *clientv3.Client
	etcdLease      clientv3.Lease
	leaseGrantResp *clientv3.LeaseGrantResponse
	KV             clientv3.KV
	serverGroup    string
	serverInst     *WindServer
}

func (sgm *ServerGroupManager) init(serverInst *WindServer) {
	client, err := clientv3.New(ETCDCONFIG)
	if err != nil {
		fmt.Println(err)
		return
	}
	sgm.etcdClient = client
	sgm.KV = clientv3.NewKV(sgm.etcdClient)
	sgm.serverGroup = SERVEARGROUPNAME
	sgm.serverInst = serverInst
}

func (sgm *ServerGroupManager) registerServerEtcd(serverId string, serverType int, etcdTTl int) {
	sgm.etcdLease = clientv3.NewLease(sgm.etcdClient)
	leaseGrantResp, err := sgm.etcdLease.Grant(context.TODO(), int64(etcdTTl))
	if err != nil {
		fmt.Println(err)
		return
	}
	nodeKey := "/"
	strings.Join([]string{nodeKey, sgm.serverGroup, "/servers/", string(rune(serverType)), "/", serverId}, "")
	info := sgm.serverInst.GetReportInfo()
	_, err = sgm.KV.Put(context.TODO(), nodeKey, info, clientv3.WithLease(leaseGrantResp.ID))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("update info to etcd", serverType, serverId, info)
}
