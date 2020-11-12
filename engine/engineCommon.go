package engine

type requestInterface struct{
	playerId 		string
	serverType 		int
	command  		string
	data            []byte		// protobuf binary
}

const ServerMaxConnect = 3000
