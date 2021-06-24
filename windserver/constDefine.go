package windserver

type ServerMetaInfo struct {
	Ip    string
	Port  int
	IntId int
}

type requestData struct {
	ServerType    		int
	Pid 		  		string
	Command       		string
	Data          		*byte     // protobuf的二进制
}

type ServerAlias = int

type serverType struct {
	INVALID 	ServerAlias		`name:"INVALID"`
	LOGIN 		ServerAlias		`name:"LOGIN"`
	LOGIC 		ServerAlias		`name:"LOGIC"`
}

// Enum for public use
var SERVERTYPE = &serverType{
	INVALID: 0,
	LOGIN: 1,
	LOGIC: 2,
}

func (st *serverType) GetServerTypeByName(name string) ServerAlias {
	switch name {
	case "loginSrv":
		return st.LOGIN
	case "logicSrv":
		return st.LOGIC
	default:
		return st.INVALID
	}
}




