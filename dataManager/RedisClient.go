package dataManager

import "github.com/go-redis/redis"


type WindServerRedisClient struct {
	Conn        	*redis.ClusterClient
}

func  (client *WindServerRedisClient) ClientInit(addr string,password string) {
	client.Conn = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"addr:7000", "addr:7001", "addr:7002", "addr:7003", "addr:7004", "addr:7005"},
		// To route commands by latency or randomly, enable one of the following.
		//RouteByLatency: true,
		//RouteRandomly: true,
	})

}