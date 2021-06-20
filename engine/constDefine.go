package engine

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

type SERVERTYPE int32

const (
	INVALID     	SERVERTYPE = 0
	LOGIN      		SERVERTYPE = 2
	LOGIC      		SERVERTYPE = 3
)
