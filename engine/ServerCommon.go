package engine

type RequestMessage struct{
	playerId 		string
	serverType 		int
	sid 			string
	command  		string
	data            []byte		// protobuf binary
}

const SERVERMAXCONNECT = 3000
const USEREDISCLUSTER = true

const (
	LOGIN	int = iota    // 开始生成枚举值, 默认为0
	GATEWAY
	LOGIC
	DB
)

func ServerTypeToName(serverType int) string {
	switch serverType {
	case LOGIN:
		return "LOGIN"
	case GATEWAY:
		return "GATEWAY"
	case LOGIC:
		return "LOGIC"
	case DB:
		return "DB"
	default:
		return "NONE"
	}
}
