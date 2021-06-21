package windserver

import (
	"context"
	"encoding/json"
	mvccpb2 "github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"strconv"
	"strings"
)


type ServerGroupManagerBasic struct {
	etcdAddr			string
	etcdGroup			string
	useGrpcProxy		bool
	etcdConfig 			clientv3.Config
	etcdClient     		*clientv3.Client
	etcdLease      		clientv3.Lease
	leaseGrantResp 		*clientv3.LeaseGrantResponse
	serverInst     		*windServer

	etcdEvent 			chan *clientv3.Event
	watchTypes			map[int]bool
	etcdWatch        	[]clientv3.WatchChan
	onlineServers		map[int]map[string]ServerMetaInfo      // server
}

func NewServerGroupManagerBasic(config clientv3.Config) *ServerGroupManagerBasic{
	return &ServerGroupManagerBasic{etcdConfig: config}
}

func (sgm *ServerGroupManagerBasic) SetUp(serverInst *windServer) {
	client, err := clientv3.New(ETCDCONFIG)
	if err != nil {
		println(err)
		return
	}
	sgm.etcdClient = client
	sgm.serverInst = serverInst
}

func (sgm *ServerGroupManagerBasic) StartService(ctx context.Context) {
	sgm.registerServerEtcd(sgm.serverInst.GetServerId(),sgm.serverInst.GetServerType(), EtcdTTl)
}

func (sgm *ServerGroupManagerBasic) registerServerEtcd(serverId string, serverType int, etcdTTl int) {
	sgm.etcdLease = clientv3.NewLease(sgm.etcdClient)
	leaseGrantResp, err := sgm.etcdLease.Grant(context.TODO(), int64(etcdTTl))
	if err != nil {
		println("update server info to etcd error:", err)
		return
	}
	var nodeKey = "/" + sgm.etcdGroup + "/servers/" + string(rune(serverType)) + "/" + serverId
	info := sgm.serverInst.GetReportInfo()
	_, err = sgm.etcdClient.KV.Put(context.TODO(), nodeKey, info, clientv3.WithLease(leaseGrantResp.ID))
	if err != nil {
		println("update server info to etcd error:", err)
		return
	}
	println("update info to etcd", serverType, serverId, info)
}

func (sgm *ServerGroupManagerBasic) AddWatch(lst []int) {
	var le = len(lst)
	for idx:=0; idx<le; idx++ {
		var serverType = lst[idx]
		sgm.watchTypes[serverType] = true
	}
}

func  (sgm *ServerGroupManagerBasic) CloseWatch()  {
	err := sgm.etcdClient.Watcher.Close()
	if err!= nil {
		println("watcher close error", err)
	}
}

func  (sgm *ServerGroupManagerBasic) WatchServers(ctx context.Context)  {
	var prefix = "/" + sgm.etcdGroup + "/servers/"
	for serverType := range sgm.watchTypes {
		sgm.onlineServers[serverType] = make(map[string]ServerMetaInfo)
		var node = prefix + string(rune(serverType)) + "/"
		var watchChan = sgm.etcdClient.Watcher.Watch(ctx, node)
		sgm.etcdWatch = append(sgm.etcdWatch, watchChan)
		go sgm.ProcessOneWatchChan(ctx, watchChan)
	}
	sgm.UpdateWatchServers()
}

func  (sgm *ServerGroupManagerBasic) ProcessOneWatchChan(ctx context.Context, watchRespChan clientv3.WatchChan)  {
	for !sgm.serverInst.serverExited {
		select {
		case <-ctx.Done():
				return
		case watchResp := <-watchRespChan:
			for _,event := range watchResp.Events {
				sgm.etcdEvent <- event
			}
		}
	}
}

func  (sgm *ServerGroupManagerBasic) UpdateWatchServers()  {
	for serverType := range sgm.watchTypes {
		sgm.UpdateServersByType(serverType)
	}
}

func  (sgm *ServerGroupManagerBasic) UpdateServersByType(serverType int)  {
	curServer := sgm.onlineServers[serverType]
	for sid,info := range curServer {
		var jsonInfo,err = json.Marshal(info)
		if err != nil {
			println("UpdateServersByType.sid:",sid," info:",jsonInfo)
		} else {
			println()
		}
	}
}

func (sgm *ServerGroupManagerBasic) ProcessEtcdEvents(ctx context.Context) {
	for !sgm.serverInst.serverExited {
		select {
		case <-ctx.Done():
			return
		case e := <- sgm.etcdEvent:
			sgm.ProcessOneEtcdEvent(e)
		}
	}
}

func (sgm *ServerGroupManagerBasic) ProcessOneEtcdEvent(event *clientv3.Event) {
	var param = strings.Split(string(event.Kv.Key), "/")
	serverType, err := strconv.Atoi(param[len(param) -2])
	if err != nil {
		println(err)
		return
	}
	curServers,ok := sgm.onlineServers[serverType]
	if !ok {
		return
	}
	var sid = param[len(param)-1]
	_,has := curServers[sid]
	switch event.Type {
	case mvccpb2.Event_EventType(mvccpb.PUT):
		var value = event.Kv.Value
		var dat map[string]interface{}
		err := json.Unmarshal(value, &dat)
		if err != nil {
			println("json.Unmarshal value:",value, " fail")
			return
		}
		var info = ServerMetaInfo{}
		info.ip = dat["ip"].(string)
		info.port = dat["port"].(int)
		info.intId = dat["intId"].(int)
		curServers[sid] = info
		sgm.onServerAdd(sid)
	case mvccpb2.Event_EventType(mvccpb.DELETE):
		if has {
			delete(curServers,sid)
			sgm.onServerDelete(sid)
		}
	}
}

func (sgm *ServerGroupManagerBasic) onServerDelete(sid string) {

}

func (sgm *ServerGroupManagerBasic) onServerAdd(sid string) {

}

func (sgm *ServerGroupManagerBasic) CheckServerOnline(sid string, serverType int) bool {
	srvs,has := sgm.onlineServers[serverType]
	if has {
		_,in := srvs[sid]
		if in {
			return true
		}
	}
	return false
}

func (sgm *ServerGroupManagerBasic) cleanEtcd(ctx context.Context) {
	var serverType = sgm.serverInst.GetServerType()
	var serverId = sgm.serverInst.GetServerId()
	nodeKey := "/" + sgm.etcdGroup + "/servers/" + string(rune(serverType)) + "/" + serverId
	_,err := sgm.etcdClient.KV.Delete(ctx,nodeKey)
	if err!=nil {
		println("error in clean Etcd")
	}
}