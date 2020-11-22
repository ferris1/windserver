package dataManager

const RedisDbName = "WindServer"

func MakeRedisKey(prefix string,targetHash string,target string) string {
	sep1 := "{"
	sep2 := "}"
	return RedisDbName + prefix + sep1 + targetHash +  sep2 + target
}

func GetPlayerServer(playerId string, serverType int) string {
	return ""
}

