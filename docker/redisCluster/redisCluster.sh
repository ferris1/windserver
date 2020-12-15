#!/bin/sh

LOCAL_IP=$2

# 这里写死
SERVERNAME="lx33"
IMAGESNAME="lx33"

SERVER_PUBLIC_NAME="redis"

errorexit()
{
	echo "Usage: $0 [create|remove|start|stop] (local_ip)"
	exit
}

create_service()
{
	echo "---------- create $SERVER_PUBLIC_NAME"
	if [ -z $LOCAL_IP ]; then
		LOCAL_IP=$(ip route get 8.8.8.8 | head -1 | awk '{print $7}')
	fi


	RedisNodePorts=(7000 7001 7002 7003 7004 7005)
	RedisClusterNodes=""
	for RedisNodePort in ${RedisNodePorts[@]}; do
		echo "---------- create redis node $RedisNodePort"
		RedisClusterNodes="$RedisClusterNodes $LOCAL_IP:$RedisNodePort"
		PreRedisCmd="echo cluster-announce-ip $LOCAL_IP >> /etc/redis/redis.conf; echo replica-announce-ip $LOCAL_IP >> /etc/redis/redis.conf; echo port $RedisNodePort >> /etc/redis/redis.conf"
		docker run --name $SERVERNAME-redis-$RedisNodePort --network host -d ${IMAGESNAME}_redis_cluster \
			sh -c "$PreRedisCmd; redis-server /etc/redis/redis.conf"
	done
	echo "---------- create redis cluster on $SERVERNAME-redis-7000"
	docker exec -d $(docker ps -a -q -f name=$SERVERNAME-redis-${RedisNodePorts[0]}) sh -c "echo yes| redis-cli --cluster create $RedisClusterNodes --cluster-replicas 1"
	echo "---------- run redis cluster finished"
}
start_service()
{
	echo "---------- start $SERVER_PUBLIC_NAME"
	docker start $(docker ps -a -q -f name=$SERVERNAME-redis)
}

stop_service()
{
	echo "---------- stop $SERVER_PUBLIC_NAME"
	docker stop $(docker ps -a -q -f name=$SERVERNAME-redis)
}

remove_service()
{
	stop_service
	echo "---------- remove $SERVER_PUBLIC_NAME"
	docker rm -v $(docker ps -a -q -f name=$SERVERNAME-redis)
}

case $1 in
	'create')
		create_service
	;;
	'remove')
		remove_service
	;;
	'start')
		start_service
	;;
	'stop')
		stop_service
	;;
	*)
		errorexit
	;;
esac
