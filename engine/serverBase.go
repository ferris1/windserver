package engine

import "github.com/ferris1/windserver/utilize"

// engine更多做的是单个进程中服务消息分发
type ServerBase struct {
	ServerName 				string
	ServerId   				string
	ServerType 				int
	ConnectCount			int
	TotalConnectCount		int
	// player 请求信息发送
	requestQueue 			utilize.EsQueue
	// 消息回调函数注册 应该是个回调函数
	commandMap				map[string]string
	// 服务器组管理
	serverGroupMgr 			ServerGroupManager
	// 连接管理
	connMgr					ConnManager
	// 消息中间件
	natsClient				NatsClient
}

//启动之前的设置
func (s *ServerBase) SetUp() {
	// 注册服务器信息,监听服务,启动心跳
	// 连接消息中间件,报告服务器压力

}

// RPC框架
func (s *ServerBase) Register() {
	// 回调信息注册

}

// 启动服务器
func (s *ServerBase) StartService() {
	// 到etcd中注册服务器信息
	// 启动消息处理线程
}

// 退出服务器
func (s *ServerBase) ExitService() {

}
