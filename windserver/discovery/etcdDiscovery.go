package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/ferris1/windserver/windserver"
	"go.etcd.io/etcd/clientv3"
	"os"
	"strconv"
	"strings"
)


type EtcdDiscovery struct {
	etcdAddr       		string
	etcdGroup      		string
	useGrpcProxy   		bool
	etcdConfig     		clientv3.Config
	etcdClient     		*clientv3.Client
	kv 					clientv3.KV
	watcher 			clientv3.Watcher
	etcdLease      		clientv3.Lease
	leaseGrantResp 		*clientv3.LeaseGrantResponse
	etcdLeaseTTl 		int
	etcdEvent 			chan clientv3.Event
	watchTypes			map[int]bool
	etcdWatch        	[]clientv3.WatchChan
	onlineServers		map[int]map[string]windserver.ServerMetaInfo // server

	options   			Options
}

func NewDiscovery(opts ...Option) Discovery {
	etcd := &EtcdDiscovery{
		options:  Options{},
		register: make(map[string]uint64),
		leases:   make(map[string]clientv3.LeaseID),
	}
	return etcd
}

func getEnvConfig(opts ...Option) []Option {
	username, password := os.Getenv("ETCD_USERNAME"), os.Getenv("ETCD_PASSWORD")
	if len(username) > 0 && len(password) > 0 {
		opts = append(opts, Auth(username, password))
	}
	address := os.Getenv("MICRO_REGISTRY_ADDRESS")
	if len(address) > 0 {
		opts = append(opts, Addrs(address))
	}
	return opts
}

func (sgm *EtcdDiscovery) SetUp() {
	client, err := clientv3.New(windserver.ETCDCONFIG)
	if err != nil {
		println(err)
		return
	}
	sgm.etcdLease = nil
	sgm.watchTypes = make(map[int]bool)
	sgm.onlineServers = make(map[int]map[string]windserver.ServerMetaInfo)
	sgm.etcdClient = client
	sgm.etcdEvent = make(chan clientv3.Event)
	sgm.kv = clientv3.NewKV(client)
	sgm.watcher = clientv3.NewWatcher(client)
}

func (sgm *EtcdDiscovery) StartService(ctx context.Context) {
	go sgm.ProcessEtcdEvents(ctx)
	sgm.WatchServers(ctx)
	sgm.registerServerEtcd(ctx,sgm.srv.GetServerId(),sgm.srv.GetServerType(), sgm.etcdLeaseTTl)
}

func (sgm *EtcdDiscovery) registerServerEtcd(ctx context.Context,serverId string, serverType int, etcdTTl int) {
	sgm.etcdLease = clientv3.NewLease(sgm.etcdClient)
	var err error
	sgm.leaseGrantResp, err = sgm.etcdLease.Grant(ctx, int64(etcdTTl))
	if err != nil {
		println("update server info to etcd error:", err)
		return
	}
	var nodeKey = "/" + sgm.etcdGroup + "/servers/" + strconv.Itoa(serverType) + "/" + serverId
	info := sgm.srv.GetReportInfo()
	_, err = sgm.kv.Put(ctx, nodeKey, info, clientv3.WithLease(sgm.leaseGrantResp.ID))
	if err != nil {
		println("update server info to etcd error:", err)
		return
	}
	println("update info to etcd", serverType, serverId, info, nodeKey)
}

func (sgm *EtcdDiscovery) AddWatch(lst []int) {
	var le = len(lst)
	for idx:=0; idx<le; idx++ {
		var serverType = lst[idx]
		sgm.watchTypes[serverType] = true
	}
}

func  (sgm *EtcdDiscovery) CloseWatch()  {
	err := sgm.watcher.Close()
	if err!= nil {
		println("watcher close error", err)
	}
}

func  (sgm *EtcdDiscovery) WatchServers(ctx context.Context)  {
	var prefix = "/" + sgm.etcdGroup + "/servers/"
	for serverType := range sgm.watchTypes {
		sgm.onlineServers[serverType] = make(map[string]windserver.ServerMetaInfo)
		serverType := serverType
		var node = prefix + strconv.Itoa(serverType) + "/"
		watchRespChan := sgm.watcher.Watch(ctx, node,clientv3.WithPrefix())
		go sgm.ProcessOneWatchChan(ctx, watchRespChan)
	}
	sgm.UpdateWatchServers()
}

func  (sgm *EtcdDiscovery) ProcessOneWatchChan(ctx context.Context, watchRespChan clientv3.WatchChan)  {
	for !sgm.srv.serverExited {
		select {
		case <-ctx.Done():
				return
		case watchResp := <-watchRespChan:
			for _,event := range watchResp.Events {
				println("the event ", string(event.Type), string(event.Kv.Key), string(event.Kv.Value))
				sgm.etcdEvent <- *event
			}
		}
	}
}

func  (sgm *EtcdDiscovery) UpdateWatchServers()  {
	println("Update Watch Servers")
	for serverType := range sgm.watchTypes {
		sgm.UpdateServersByType(serverType)
	}
}

func  (sgm *EtcdDiscovery) UpdateServersByType(serverType int)  {
	curServer := sgm.onlineServers[serverType]
	for sid,info := range curServer {
		var jsonInfo,err = json.Marshal(info)
		if err != nil {
			println("UpdateServersByType.sid:",sid," info:",jsonInfo)
		} else {
			println("err when update server")
		}
	}
}

func (sgm *EtcdDiscovery) ProcessEtcdEvents(ctx context.Context) {
	for !sgm.srv.serverExited {
		select {
		case <-ctx.Done():
			return
		case e := <- sgm.etcdEvent:
			sgm.ProcessOneEtcdEvent(e)
		}
	}
}

func (sgm *EtcdDiscovery) ProcessOneEtcdEvent(event clientv3.Event) {
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
		var info = windserver.ServerMetaInfo{}
		info.Ip = dat["Ip"].(string)
		info.Port = int(dat["Port"].(float64))
		info.IntId = int(dat["IntId"].(float64))
		curServers[sid] = info
		sgm.onServerAdd(sid, string(value))
	case mvccpb.DELETE:
		if has {
			delete(curServers,sid)
			sgm.onServerDelete(sid)
		}
	}
}

func (sgm *EtcdDiscovery) onServerDelete(sid string) {
	println("onServerDelete:",sid)
}

func (sgm *EtcdDiscovery) onServerAdd(sid string, info string) {
	println("onServerAdd:",sid, " info:",info)
}

func (sgm *EtcdDiscovery) CheckServerOnline(sid string, serverType int) bool {
	srvs,has := sgm.onlineServers[serverType]
	if has {
		_,in := srvs[sid]
		if in {
			return true
		}
	}
	return false
}

func (sgm *EtcdDiscovery) CleanEtcd(ctx context.Context) {
	var serverType = sgm.srv.GetServerType()
	var serverId = sgm.srv.GetServerId()
	nodeKey := "/" + sgm.etcdGroup + "/servers/" + strconv.Itoa(serverType) + "/" + serverId
	_,err := sgm.kv.Delete(ctx, nodeKey)
	if err!=nil {
		println("error in clean Etcd")
	}
}

func (sgm *EtcdDiscovery) EtcdTick(ctx context.Context) {
	if sgm.srv.serverExited  {
		return
	}
	if sgm.etcdLease == nil || sgm.etcdLeaseTTl == 0 {
		return
	}
	if keepRespChan, err := sgm.etcdLease.KeepAliveOnce(ctx, sgm.leaseGrantResp.ID); err != nil {
		fmt.Println(err)
		sgm.etcdLease = nil
		return
	} else {
		if keepRespChan!=nil {
			println("etcd Keep Alive success")
		}
	}

}