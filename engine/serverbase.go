package engine

type ServerBase struct {
	ServerName string
	ServerId   string
}

//启动之前的设置
func (s *ServerBase) SetUp() {

}

// 注册消息
func (s *ServerBase) Register() {

}

// 启动服务器
func (s *ServerBase) StartService() {

}

// 退出服务器
func (s *ServerBase) ExitService() {

}
