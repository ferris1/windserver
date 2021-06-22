package windserver

import (
	"go.etcd.io/etcd/clientv3"
)

const RedisClusterIp = "127.0.0.1"

const SERVEARGROUPNAME = "windServer"

var (
	ETCDCONFIG = clientv3.Config{
		Endpoints: []string{"192.168.0.106:2379", "192.168.0.106:2479", "192.168.0.106:2579"},
	}
)

const EtcdTTl = 30

const SERVERMAXCONNECT = 1000

const REQUSETQUEUELEN = 100
