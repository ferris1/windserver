package main

func main() {
	var server *LogicSrv
	server = NewLogicSrv("logicSrv")
	server.Run()
	defer server.Stop()
}
