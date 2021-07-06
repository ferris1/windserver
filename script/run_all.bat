start "logicSrv" cmd /c "go run ..\service\logicSrv\logicSrv.go 50100 & pause"

start "logicSrv1" cmd /c "go run ..\service\logicSrv\logicSrv.go 50200 & pause"

start "logicSrv2" cmd /c "go run ..\service\logicSrv\logicSrv.go 50300 & pause"

pause