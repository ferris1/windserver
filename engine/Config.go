package engine

import (
	"go.etcd.io/etcd/clientv3"
	"time"
)

const RedisClusterIp = "127.0.0.1"
const SERVEARGROUP = "WindServer"
var (
	RedisClusterConf = []string{RedisClusterIp +":7000", RedisClusterIp+":7001", RedisClusterIp + ":7002",
		RedisClusterIp+ ":7003", RedisClusterIp+":7004", RedisClusterIp +":7005"}

	ETCDCONFIG = clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	}


)


