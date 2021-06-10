package engine

import (
	"encoding/json"
)

type WindServerBase struct {
	serverIp				string
	serverPort				int
	serverName 				string
	serverId   				string
	intId					int
	serverType 				int
	connectCount			int
	totalConnectCount		int
	serverExited			bool
	// player 请求信息发送
	requestChannel 			chan RequestMessage
	// 消息回调函数注册 应该是个回调函数
	commandMap				map[string]string
	// 服务器组管理
	serverGroupMgr 			ServerGroupManager
	// 客戶端连接管理
	connMgr					ConnManager
}

//启动之前的设置
func (s *WindServerBase) SetUp() {
	// 注册服务器信息,监听服务,启动心跳
	// 连接消息中间件,报告服务器压力
	// 数据初始化
	s.totalConnectCount = SERVERMAXCONNECT
}

// RPC框架  这个实现上好像比较麻烦
// 事件注册:如网络事件注册
func (s *WindServerBase) Register() {

}

// 启动服务器,一些服务线程将在这里启动,一些定时任务在这里驱动
func (s *WindServerBase) StartService() {
	// 到etcd中注册服务器信息
	// 启动消息处理线程
	s.serverGroupMgr.registerServerEtcd(s.serverId, s.serverType, EtcdTTl)

}

// 退出服务器
func (s *WindServerBase) ExitService() {

}

func (s *WindServerBase) StartUp() {
	s.SetUp()
	s.Register()
	s.StartService()
}

func (s *WindServerBase) GetReportInfo() string {
	var info map[string]string
	info = make(map[string]string)
	info["Ip"] = s.serverIp
	info["Port"] = string(rune(s.serverPort))
	info["IntId"] = string(rune(s.intId))
	res, err := json.Marshal(info)
	if err != nil {
		return "{}"
	}
	return string(res)
}

