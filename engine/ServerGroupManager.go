package engine

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"strings"
)

type ServerGroupManagerBasic struct {
	etcdAddr			string
	etcdGroup			string
	useGrpcProxy		bool
	etcdClient     		*clientv3.Client
	etcdLease      		clientv3.Lease
	leaseGrantResp 		*clientv3.LeaseGrantResponse
	serverInst     		*windServer

	watchTypes			map[int]bool
	watchServers        []int
	etcdEvent			chan string
	onlineServers		map[int][]ServerMetaInfo      // server
}

func NewServerGroupManagerBasic(config EtcdConfig) *ServerGroupManagerBasic{
	return &ServerGroupManagerBasic{etcdAddr: config.EtcdAddr,
		etcdGroup: config.EtcdGroup, useGrpcProxy: config.UseGrpcProxy}
}

func (sgm *ServerGroupManagerBasic) SetUp(serverInst *windServer) {
	client, err := clientv3.New(ETCDCONFIG)
	if err != nil {
		fmt.Println(err)
		return
	}
	sgm.etcdClient = client
	sgm.serverInst = serverInst
}

func (sgm *ServerGroupManagerBasic) StartService() {

}

func (sgm *ServerGroupManagerBasic) ProcessEtcdEvents() {
	for !sgm.serverInst.serverExited {
		select {
		case e := <- sgm.etcdEvent:
			fmt.Println("event:",e)
		}
	}
}

func (sgm *ServerGroupManagerBasic) ProcessOneEtcdEvents() {

}

func (sgm *ServerGroupManagerBasic) registerServerEtcd(serverId string, serverType int, etcdTTl int) {
	sgm.etcdLease = clientv3.NewLease(sgm.etcdClient)
	leaseGrantResp, err := sgm.etcdLease.Grant(context.TODO(), int64(etcdTTl))
	if err != nil {
		fmt.Println(err)
		return
	}
	nodeKey := "/"
	strings.Join([]string{nodeKey, sgm.etcdGroup, "/servers/", string(rune(serverType)), "/", serverId}, "")
	info := sgm.serverInst.GetReportInfo()
	_, err = sgm.etcdClient.Put(context.TODO(), nodeKey, info, clientv3.WithLease(leaseGrantResp.ID))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("update info to etcd", serverType, serverId, info)
}
