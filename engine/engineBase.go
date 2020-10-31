package engine

import "github.com/ferris1/windserver/utilize"

type EngineBase struct {
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
func (s *EngineBase) SetUp() {
	// 注册服务器信息,监听服务,启动心跳
	// 连接消息中间件,报告服务器压力

}

// RPC框架  这个实现上好像比较麻烦
func (s *EngineBase) Register() {


}

// 启动服务器
func (s *EngineBase) StartService() {
	// 到etcd中注册服务器信息
	// 启动消息处理线程
}

// 退出服务器
func (s *EngineBase) ExitService() {

}