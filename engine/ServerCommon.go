package engine

type RequestMessage struct{
	playerId 		string
	serverType 		int
	command  		string
	data            []byte		// protobuf binary
}

const ServerMaxConnect = 3000
const UseRedisCluster = true
