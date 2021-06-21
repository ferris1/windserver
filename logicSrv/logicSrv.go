package main

import (
	"github.com/ferris1/windserver/windserver"
)

type LogicSrv struct {
	windserver.WindServer
}

func  NewLogicSrv(name string)  *LogicSrv {
	return &LogicSrv{windserver.NewWindServer(name)}
}

func (s *LogicSrv) SetUp() {
	s.WindServer.SetUp()
	s.AddWatchServers([]int{windserver.SERVERTYPE.LOGIC})
	println("server has SetUp....")
}

// RPC框架  这个实现上好像比较麻烦
// 事件注册:如网络事件注册
func (s *LogicSrv) Register() {
	s.WindServer.Register()
	println("logic server has Register....")
}

// 启动服务器,一些服务线程将在这里启动,一些定时任务在这里驱动
func (s *LogicSrv) StartService() {
	s.WindServer.StartService()
	println("logic server start service....")
}

// 退出服务器
func (s *LogicSrv) Stop() {
	s.WindServer.Stop()
	println("logic server exit service....")
}

func (s *LogicSrv) Run() {
	s.SetUp()
	s.Register()
	s.StartService()
}

