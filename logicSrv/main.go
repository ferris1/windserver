package main

func main() {
	var server *LogicSrv
	server = NewLogicSrv("logic")
	server.Run()
	defer server.Stop()
}
