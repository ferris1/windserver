package main

import (
	"github.com/ferris1/windserver/engine"
)

type LogicSrv struct {
	base *engine.WindServer
}

func  NewLogicSrv(name string)  *LogicSrv {
	return &LogicSrv{base: engine.New(name)}
}

func (s *LogicSrv) SetUp() {
	s.base.SetUp()
	println("server has SetUp....")
}

// RPC框架  这个实现上好像比较麻烦
// 事件注册:如网络事件注册
func (s *LogicSrv) Register() {
	s.base.Register()
	println("logic server has Register....")
}

// 启动服务器,一些服务线程将在这里启动,一些定时任务在这里驱动
func (s *LogicSrv) StartService() {
	s.base.StartService()
	println("logic server start service....")
}

// 退出服务器
func (s *LogicSrv) ExitService() {
	s.base.ExitService()
	println("logic server exit service....")
}

func (s *LogicSrv) Run() {
	s.SetUp()
	s.Register()
	s.StartService()
}

