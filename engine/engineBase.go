package engine

type windServerBase struct {
	serverName 				string
	serverId   				string
	serverType 				int
	connectCount			int
	totalConnectCount		int
	serverExited			bool
	// player 请求信息发送
	requestChannel 			chan requestInterface
	// 消息回调函数注册 应该是个回调函数
	commandMap				map[string]string
	// 服务器组管理
	serverGroupMgr 			ServerGroupManager
	// 客戶端连接管理
	connMgr					ConnManager
	// 消息中间件
	nsqClient				nsqClient
}

//启动之前的设置
func (s *windServerBase) SetUp() {
	// 注册服务器信息,监听服务,启动心跳
	// 连接消息中间件,报告服务器压力
	// 数据初始化
	s.totalConnectCount = ServerMaxConnect
	s.requestChannel = make(chan requestInterface,s.totalConnectCount)
}

// RPC框架  这个实现上好像比较麻烦
// 事件注册:如网络事件注册
func (s *windServerBase) Register() {


}

// 启动服务器,一些服务线程将在这里启动,一些定时任务在这里驱动
func (s *windServerBase) StartService() {
	// 到etcd中注册服务器信息
	// 启动消息处理线程
	go s.handleRequestQueue()
}

// 退出服务器
func (s *windServerBase) ExitService() {

}

func (s *windServerBase) handleRequestQueue() {
	for s.serverExited {

	}
}

