package engine

import (
	"go.etcd.io/etcd/clientv3"
)

const RedisClusterIp = "127.0.0.1"

const SERVEARGROUPNAME = "WindServer"

var (
	ETCDCONFIG = clientv3.Config{
		Endpoints: []string{"localhost:2379", "localhost:22379", "localhost:32379"},
	}
)

const EtcdTTl = 30

const SERVERMAXCONNECT = 1000
