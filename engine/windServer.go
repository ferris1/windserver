package engine

import (
	"context"
	"encoding/json"
	"github.com/ferris1/windserver/engine/until/signals"
)

type WindServer struct {
	serverIp          string
	serverPort        int
	serverName        string
	serverId          string
	intId             int
	serverType        int
	connectCount      int
	totalConnectCount int
	serverExited      bool
	// 消息回调函数注册 应该是个回调函数
	commandMap map[string]string
	// 服务器组管理
	serverGroupMgr ServerGroupManager
	// 客戶端连接管理
	connMgr ConnManager
}

func  New(name string)  *WindServer{
	return &WindServer{serverName: name}
}

//启动之前的设置
func (s *WindServer) SetUp() {
	// 注册服务器信息,监听服务,启动心跳
	// 连接消息中间件,报告服务器压力
	// 数据初始化

	s.totalConnectCount = SERVERMAXCONNECT
	println("wind server base has SetUp....")
}

// RPC框架  这个实现上好像比较麻烦
// 事件注册:如网络事件注册
func (s *WindServer) Register() {
	println("wind server base has Register....")
}

// 启动服务器,一些服务线程将在这里启动,一些定时任务在这里驱动
func (s *WindServer) StartService() {
	// 到etcd中注册服务器信息
	// 启动消息处理线程
	// s.serverGroupMgr.registerServerEtcd(s.serverId, s.serverType, EtcdTTl)
	// 主线程处理循环
	println("wind server base running... ")
	ctx := signals.NewSigKillContext()
	go s.ProcessMessageQueue(ctx)

}

// 退出服务器
func (s *WindServer) ExitService() {
	println("wind server base has ExitService....")
}

func (s *WindServer) Run() {
	s.SetUp()
	s.Register()
	s.StartService()
}

func (s *WindServer) GetReportInfo() string {
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

func (s *WindServer) ProcessMessageQueue(ctx context.Context)  {
	for !s.serverExited {
		select {
		case <-ctx.Done():
			return
		default:

		}
	}
}

// func main() {
// 	var server WindServer
// 	server.StartUp()
// }
