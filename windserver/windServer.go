package windserver

import (
	"context"
	"encoding/json"
	"github.com/ferris1/windserver/windserver/until/signals"
)

type WindServer interface {
	// MemberList lists the current cluster membership.
	SetUp()
	Register()
	StartService()
	Run()
	Stop()
	AddWatchServers(lst []int)
}

// 引擎主要功能：
// 1.定义服务器框架的流程
// 2.初始化网络模块
// 3.配合组件server_group_manager和switcher_client，建立服务器组的信息同步、数据收发和负载监控
// 4.消息注册和向上转发等
type windServer struct {
	serverIp          	string
	serverPort        	int
	serverName        	string
	serverId          	string
	intId             	int
	serverType        	int
	connectCount      	int
	totalConnectCount 	int
	serverExited      	bool
	requestQueue      	chan requestData
	// 消息回调函数注册 应该是个回调函数
	commandMap 			map[string]string
	// 客戶端连接管理
	connMgr 			*ConnManager
	serverGroupMgr 		*ServerGroupManagerBasic
}

func NewWindServer(name string)  WindServer {
	return &windServer{serverName: name, serverId: "12345678", serverType: 1}
}

//启动之前的设置
func (s *windServer) SetUp() {
	// 注册服务器信息,监听服务,启动心跳
	// 连接消息中间件,报告服务器压力
	// 数据初始化
	s.serverExited = false
	s.totalConnectCount = SERVERMAXCONNECT
	s.serverGroupMgr = NewServerGroupManagerBasic(ETCDCONFIG, "WindServer")
	s.serverGroupMgr.SetUp(s)
	println("wind server base has SetUp....")
}

// RPC框架  这个实现上好像比较麻烦
// 事件注册:如网络事件注册
func (s *windServer) Register() {
	println("wind server base has Register....")
}

// 启动服务器,一些服务线程将在这里启动,一些定时任务在这里驱动
func (s *windServer) StartService() {
	// 到etcd中注册服务器信息
	// 启动消息处理线程
	// s.serverGroupMgr.registerServerEtcd(s.serverId, s.serverType, EtcdTTl)
	// 主线程处理循环

	println("wind server base running... ")
	ctx := signals.NewSigKillContext()
	s.serverGroupMgr.StartService(ctx)
	go s.ProcessMessageQueue(ctx)
	<-ctx.Done()
}

func (s *windServer) Run() {
	s.SetUp()
	s.Register()
	s.StartService()
}

// 退出服务器
func (s *windServer) Stop() {
	println("wind server base has ExitService....")
}

func (s *windServer) GetServerId() string {
	return s.serverId
}

func (s *windServer) GetServerType() int {
	return s.serverType
}

func (s *windServer) GetReportInfo() string {
	var info = &ServerMetaInfo{}
	info.Ip = s.serverId
	info.Port = s.serverPort
	info.IntId = s.intId
	res, err := json.Marshal(info)
	if err != nil {
		println("GetReportInfo.err:",err)
		return "{}"
	}
	return string(res)
}

// 发送消息时是需要保证时序
// 可以在增加一个不需要保证顺序的接口
func (s *windServer) ProcessMessageQueue(ctx context.Context)  {
	for !s.serverExited {
		select {
		case <-ctx.Done():
			return
		case req := <- s.requestQueue:
			var sid,ok = s.GetCommandDestSid(req.Pid, req.ServerType, req.Command)
			if ok {
				s.SendDataByMessageServer(req.Pid, sid, req.ServerType, req.Command, req.Data)
			}
		}
	}
}

func (s *windServer) SendDataByMessageServer(pid string, sid string, serverType int, command string, data *byte)  bool {
	return true
}

func (s *windServer) GetCommandDestSid(pid string, serverType int, command string)  (string,bool) {
	return "HAHA",true
}

func (s *windServer) AddWatchServers(lst []int) {
	s.serverGroupMgr.AddWatch(lst)
}

// func main() {
// 	var server windServer
// 	server.StartUp()
// }
