package windserver

type ServerMetaInfo struct {
	ip    				string
	port 				int
	intId				int
}

type requestData struct {
	ServerType    		int
	Pid 		  		string
	Command       		string
	Data          		*byte     // protobuf的二进制
}

type ServerAlias = int

type serverType struct {
	INVALID 	ServerAlias
	LOGIN 		ServerAlias
	LOGIC 		ServerAlias
}

// Enum for public use
var SERVERTYPE = &serverType{
	INVALID: 0,
	LOGIN: 1,
	LOGIC: 2,
}

