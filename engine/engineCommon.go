package engine

type requestMessage struct{
	playerId 		string
	serverType 		int
	command  		string
	data            []byte		// protobuf binary
}

const ServerMaxConnect = 3000
