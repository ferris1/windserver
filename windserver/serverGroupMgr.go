package windserver

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3"
	"strconv"
	"strings"
)


type ServerGroupManagerBasic struct {
	etcdAddr       string
	etcdGroup      string
	useGrpcProxy   bool
	etcdConfig     clientv3.Config
	etcdClient     *clientv3.Client
	etcdLease      clientv3.Lease
	leaseGrantResp *clientv3.LeaseGrantResponse
	srv            *windServer

	etcdEvent 			chan *clientv3.Event
	watchTypes			map[int]bool
	etcdWatch        	[]clientv3.WatchChan
	onlineServers		map[int]map[string]ServerMetaInfo      // server
}

func NewServerGroupManagerBasic(config clientv3.Config, etcdGroup string) *ServerGroupManagerBasic{
	return &ServerGroupManagerBasic{etcdConfig: config,etcdGroup: etcdGroup}
}

func (sgm *ServerGroupManagerBasic) SetUp(serverInst *windServer) {
	client, err := clientv3.New(ETCDCONFIG)
	if err != nil {
		println(err)
		return
	}
	sgm.watchTypes = make(map[int]bool)
	sgm.onlineServers = make(map[int]map[string]ServerMetaInfo)
	sgm.etcdClient = client
	sgm.srv = serverInst
}

func (sgm *ServerGroupManagerBasic) StartService(ctx context.Context) {
	sgm.registerServerEtcd(ctx,sgm.srv.GetServerId(),sgm.srv.GetServerType(), EtcdTTl)
	sgm.WatchServers(ctx)
}

func (sgm *ServerGroupManagerBasic) registerServerEtcd(ctx context.Context,serverId string, serverType int, etcdTTl int) {
	sgm.etcdLease = clientv3.NewLease(sgm.etcdClient)
	leaseGrantResp, err := sgm.etcdLease.Grant(ctx, int64(etcdTTl))
	if err != nil {
		println("update server info to etcd error:", err)
		return
	}
	var nodeKey = "/" + sgm.etcdGroup + "/servers/" + strconv.Itoa(serverType) + "/" + serverId
	info := sgm.srv.GetReportInfo()
	_, err = sgm.etcdClient.KV.Put(ctx, nodeKey, info, clientv3.WithLease(leaseGrantResp.ID))
	if err != nil {
		println("update server info to etcd error:", err)
		return
	}
	println("update info to etcd", serverType, serverId, info, nodeKey)
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
		var node = prefix + strconv.Itoa(serverType) + "/"
		println("watch server prefix:",node)
		var watchChan = sgm.etcdClient.Watcher.Watch(ctx, node)
		sgm.etcdWatch = append(sgm.etcdWatch, watchChan)
		go sgm.ProcessOneWatchChan(ctx, watchChan)
	}
	sgm.UpdateWatchServers()
}

func  (sgm *ServerGroupManagerBasic) ProcessOneWatchChan(ctx context.Context, watchRespChan clientv3.WatchChan)  {
	for !sgm.srv.serverExited {
		println("watch chan:")
		select {
		case <-ctx.Done():
				return
		case watchResp := <-watchRespChan:
			println("watchResp")
			for _,event := range watchResp.Events {
				println("the event ", event.Type, event.Kv.Key, event.Kv.Value)
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
	for !sgm.srv.serverExited {
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
	case mvccpb.PUT:
		var value = event.Kv.Value
		var dat map[string]interface{}
		err := json.Unmarshal(value, &dat)
		if err != nil {
			println("json.Unmarshal value:",value, " fail")
			return
		}
		var info = ServerMetaInfo{}
		info.Ip = dat["Ip"].(string)
		info.Port = dat["Port"].(int)
		info.IntId = dat["IntId"].(int)
		curServers[sid] = info
		sgm.onServerAdd(sid)
	case mvccpb.DELETE:
		if has {
			delete(curServers,sid)
			sgm.onServerDelete(sid)
		}
	}
}

func (sgm *ServerGroupManagerBasic) onServerDelete(sid string) {
	println("onServerDelete:",sid)
}

func (sgm *ServerGroupManagerBasic) onServerAdd(sid string) {
	println("onServerAdd:",sid)
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
	var serverType = sgm.srv.GetServerType()
	var serverId = sgm.srv.GetServerId()
	nodeKey := "/" + sgm.etcdGroup + "/servers/" + strconv.Itoa(serverType) + "/" + serverId
	_,err := sgm.etcdClient.KV.Delete(ctx,nodeKey)
	if err!=nil {
		println("error in clean Etcd")
	}
}