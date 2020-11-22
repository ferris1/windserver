package dataManager

import "github.com/go-redis/redis"


type WindServerRedisClient struct {
	Conn        	*redis.ClusterClient
}

func  (client *WindServerRedisClient) ClientInit(addrs []string,password string) {
	client.Conn = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addrs,
		// To route commands by latency or randomly, enable one of the following.
		//RouteByLatency: true,
		//RouteRandomly: true,
	})

}