#!/bin/sh

LOCAL_IP=$2

# 这里写死
SERVERNAME="windserver"
DATA_DIR_PRE=$(cd `dirname $0`;pwd)
IMAGESNAME="windserver"

SERVER_PUBLIC_NAME="etcd"

errorexit()
{
	echo "Usage: $0 [create|remove|start|stop] (local_ip)"
	exit
}


create_service()
{

	if [ -z $LOCAL_IP ] ; then
		LOCAL_IP=$(ip route get 8.8.8.8 | head -1 | awk '{print $7}')
	fi

	echo "---------- create $SERVER_PUBLIC_NAME"

	ETCD_PORTS=(2379 2479 2579)

	etcd_id=0
	cluster_addr=
	for port in ${ETCD_PORTS[@]}
	do
		etcd_id=$((etcd_id+1))
		peer_port=$((port+1))
		if [ -z $cluster_addr ]; then
			cluster_addr="etcd"$etcd_id"=http://"${LOCAL_IP}":"$peer_port
		else
			cluster_addr=$cluster_addr",etcd"$etcd_id"=http://"${LOCAL_IP}":"$peer_port
		fi
	done

	etcd_id=0
	for port in ${ETCD_PORTS[@]}
	do
		etcd_id=$((etcd_id+1))
		peer_port=$((port+1))
		echo "---------- create etcd node $port"

		data_dir="$DATA_DIR_PRE/etcd$etcd_id"

		if [ ! -d $data_dir ]; then
			echo "The dir $data_dir isn't exist, it's being created."
			mkdir $data_dir -p
			chmod 777 $data_dir
		fi

		docker run --name $SERVERNAME-etcd$etcd_id --network host --volume=$data_dir:/etcd-data -d quay.io/coreos/etcd:v3.4.1 /usr/local/bin/etcd \
			--data-dir=/etcd-data --name etcd$etcd_id \
			--initial-advertise-peer-urls http://${LOCAL_IP}:$peer_port --listen-peer-urls http://0.0.0.0:$peer_port \
			--advertise-client-urls http://${LOCAL_IP}:$port --listen-client-urls http://0.0.0.0:$port \
			--initial-cluster-token etcd-cluster-1 --initial-cluster-state new \
			--initial-cluster $cluster_addr
	done
}

start_service()
{
	echo "---------- start $SERVER_PUBLIC_NAME"
	docker start $(docker ps -a -q -f name=$SERVERNAME-etcd)
}

stop_service()
{
	echo "---------- stop $SERVER_PUBLIC_NAME"
	docker stop $(docker ps -a -q -f name=$SERVERNAME-etcd)
}


remove_service()
{
	stop_service
	echo "---------- remove $SERVER_PUBLIC_NAME"
	docker rm -v $(docker ps -a -q -f name=$SERVERNAME-etcd)
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
