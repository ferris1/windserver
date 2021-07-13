package discovery

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/ferris1/windserver/windserver"
	"go.etcd.io/etcd/clientv3"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	prefix = "/windserver/discovery/"
)

type etcdDiscovery struct {
	sync.RWMutex
	client 				*clientv3.Client
	onlineService 		map[string]Service
	options   			Options
	lease      			clientv3.LeaseID
	register 			string
}

func NewDiscovery(opts ...Option) Discovery {
	etcd := &etcdDiscovery{
		options:  Options{},
		onlineService: make(map[int]Service),
	}
	opts = getEnvConfig(opts...)
	_ = configure(etcd, opts...)
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

func configure(e *etcdDiscovery, opts ...Option) error {
	config := clientv3.Config{
		Endpoints: etcdEndpoints,
	}

	for _, o := range opts {
		o(&e.options)
	}

	if e.options.Timeout == 0 {
		e.options.Timeout = 5 * time.Second
	}
	config.DialTimeout = e.options.Timeout

	if e.options.Secure || e.options.TLSConfig != nil {
		tlsConfig := e.options.TLSConfig
		if tlsConfig == nil {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		config.TLS = tlsConfig
	}

	if e.options.Username != "" {
		config.Username = e.options.Username
		config.Password = e.options.Password
	}

	var cAddrs []string

	for _, address := range e.options.Addrs {
		if len(address) == 0 {
			continue
		}
		addr, port, err := net.SplitHostPort(address)
		if ae, ok := err.(*net.AddrError); ok && ae.Err == "missing port in address" {
			port = "2379"
			addr = address
			cAddrs = append(cAddrs, net.JoinHostPort(addr, port))
		} else if err == nil {
			cAddrs = append(cAddrs, net.JoinHostPort(addr, port))
		}
	}

	// if we got addrs then we'll update
	if len(cAddrs) > 0 {
		config.Endpoints = cAddrs
	}

	cli, err := clientv3.New(config)
	if err != nil {
		return err
	}
	e.client = cli
	return nil
}

func nodePath(s, id string) string {
	return path.Join(prefix, s, id)
}

func encode(s *Node) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func decode(ds []byte) *Node {
	var s *Node
	json.Unmarshal(ds, &s)
	return s
}

func (e *etcdDiscovery) SetUp(opts ...Option) {
	_ = configure(e, opts...)
}

func (e *etcdDiscovery) Options() Options {
	return e.options
}

func (e *etcdDiscovery) StartService(ctx context.Context) {
	go e.ProcessEtcdEvents(ctx)
	e.WatchServers(ctx)
}

func (e *etcdDiscovery) Register(n *Node, opts ...RegisterOption) error {
	return e.registerNode(n, opts...)
}

func (e *etcdDiscovery) registerNode(node *Node, opts ...RegisterOption) error {
	var options RegisterOptions
	var lgr *clientv3.LeaseGrantResponse
	var err error
	for _, o := range opts {
		o(&options)
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	if options.TTL.Seconds() > 0 {
		// get a lease used to expire keys since we have a ttl
		lgr, err = e.client.Grant(ctx, int64(options.TTL.Seconds()))
		if err != nil {
			return err
		}
	}
	if lgr != nil {
		_, err = e.client.Put(ctx, nodePath(node.Type, node.Id), encode(node), clientv3.WithLease(lgr.ID))
	} else {
		_, err = e.client.Put(ctx, nodePath(node.Type, node.Id), encode(node))
	}
	if err != nil {
		println("update server info to etcd error:", err)
		return err
	}
	e.Lock()
	e.lease = lgr.ID
	e.register = node.Type + node.Id
	e.Unlock()
	return nil
}

func (e *etcdDiscovery) Deregister(n *Node, opts ...DeregisterOption) error {
	e.Lock()
	e.lease = 0
	e.register = ""
	e.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	_, err := e.client.Delete(ctx, nodePath(n.Type, n.Id))
	if err != nil {
		return err
	}
	return nil
}

func (e *etcdDiscovery) GetService(name string, opts ...GetOption) (*Service, error) {

}


func (e *etcdDiscovery) AddWatch(lst []int) {
	var le = len(lst)
	for idx:=0; idx<le; idx++ {
		var serverType = lst[idx]
		e.watchTypes[serverType] = true
	}
}

func  (e *etcdDiscovery) Watch(opts ...WatchOption)  {
	var options WatchOptions
	for _, o := range opts {
		o(&options)
	}
	var prefix = "/" + e.etcdGroup + "/servers/"
	for serverType := range e.watchTypes {
		e.onlineServers[serverType] = make(map[string]windserver.ServerMetaInfo)
		serverType := serverType
		var node = prefix + strconv.Itoa(serverType) + "/"
		watchRespChan := e.watcher.Watch(ctx, node,clientv3.WithPrefix())
		go e.ProcessOneWatchChan(ctx, watchRespChan)
	}
	e.UpdateWatchServers()
}

func  (e *etcdDiscovery) ProcessOneWatchChan(ctx context.Context, watchRespChan clientv3.WatchChan)  {
	for !e.srv.serverExited {
		select {
		case <-ctx.Done():
				return
		case watchResp := <-watchRespChan:
			for _,event := range watchResp.Events {
				println("the event ", string(event.Type), string(event.Kv.Key), string(event.Kv.Value))
				e.etcdEvent <- *event
			}
		}
	}
}

func  (e *etcdDiscovery) UpdateWatchServers()  {
	println("Update Watch Servers")
	for serverType := range e.watchTypes {
		e.UpdateServersByType(serverType)
	}
}

func  (e *etcdDiscovery) UpdateServersByType(serverType int)  {
	curServer := e.onlineServers[serverType]
	for sid,info := range curServer {
		var jsonInfo,err = json.Marshal(info)
		if err != nil {
			println("UpdateServersByType.sid:",sid," info:",jsonInfo)
		} else {
			println("err when update server")
		}
	}
}

func (e *etcdDiscovery) ProcessEtcdEvents(ctx context.Context) {
	for !e.srv.serverExited {
		select {
		case <-ctx.Done():
			return
		case e := <- e.etcdEvent:
			e.ProcessOneEtcdEvent(e)
		}
	}
}

func (e *etcdDiscovery) ProcessOneEtcdEvent(event clientv3.Event) {
	var param = strings.Split(string(event.Kv.Key), "/")
	serverType, err := strconv.Atoi(param[len(param) -2])
	if err != nil {
		println(err)
		return
	}
	curServers,ok := e.onlineServers[serverType]
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
		e.onServerAdd(sid, string(value))
	case mvccpb.DELETE:
		if has {
			delete(curServers,sid)
			e.onServerDelete(sid)
		}
	}
}

func (e *etcdDiscovery) onServerDelete(sid string) {
	println("onServerDelete:",sid)
}

func (e *etcdDiscovery) onServerAdd(sid string, info string) {
	println("onServerAdd:",sid, " info:",info)
}

func (e *etcdDiscovery) CheckServerOnline(sid string, serverType int) bool {
	srvs,has := e.onlineServers[serverType]
	if has {
		_,in := srvs[sid]
		if in {
			return true
		}
	}
	return false
}

func (e *etcdDiscovery) CleanEtcd(ctx context.Context) {
	var serverType = e.srv.GetServerType()
	var serverId = e.srv.GetServerId()
	nodeKey := "/" + e.etcdGroup + "/servers/" + strconv.Itoa(serverType) + "/" + serverId
	_,err := e.kv.Delete(ctx, nodeKey)
	if err!=nil {
		println("error in clean Etcd")
	}
}

func (e *etcdDiscovery) EtcdTick(ctx context.Context) {
	if e.srv.serverExited  {
		return
	}
	if e.etcdLease == nil || e.etcdLeaseTTl == 0 {
		return
	}
	if keepRespChan, err := e.etcdLease.KeepAliveOnce(ctx, e.leaseGrantResp.ID); err != nil {
		fmt.Println(err)
		e.etcdLease = nil
		return
	} else {
		if keepRespChan!=nil {
			println("etcd Keep Alive success")
		}
	}

}