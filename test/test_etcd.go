package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

func main() {
	var (
		client *clientv3.Client
		err error
		kv clientv3.KV
		watcher clientv3.Watcher
		watchStartRevision int64
		watchRespChan <-chan clientv3.WatchResponse
		watchResp clientv3.WatchResponse
		event *clientv3.Event
	)
	ETCDCONFIG := clientv3.Config{
		Endpoints: []string{"192.168.0.106:2379", "192.168.0.106:2479", "192.168.0.106:2579"},
		DialTimeout: time.Second,
	}
	client, err = clientv3.New(ETCDCONFIG)
	if err != nil {
		fmt.Println("connect failed err : ", err)
		return
	}
	defer client.Close()
	ctx := context.Background()
	// KV
	kv = clientv3.NewKV(client)

	// 模拟etcd中KV的变化
	go func(ctx context.Context) {
		for {
			_, _ = kv.Put(ctx, "/cron/jobs/job7", "i am job7")

			_, _ = kv.Delete(ctx, "/cron/jobs/job7")

			time.Sleep(1 * time.Second)
		}
	}(ctx)

	// 先GET到当前的值，并监听后续变化
	if getResp, err := kv.Get(ctx, "/cron/jobs/job7"); err != nil {
		fmt.Println(err)
		return
	} else {
		for _,x := range getResp.Kvs {
			println("x:",string(x.Key),string(x.Value))
		}
		watchStartRevision = getResp.Header.Revision + 1
	}


	// 当前etcd集群事务ID, 单调递增的（监听/cron/jobs/job7后续的变化,也就是通过监听版本变化）


	// 创建一个watcher(监听器)
	watcher = clientv3.NewWatcher(client)

	// 启动监听
	fmt.Println("从该版本向后监听:", watchStartRevision)

	ctx, cancelFunc := context.WithCancel(context.TODO())
	//5秒钟后取消
	time.AfterFunc(10 * time.Second, func() {
		cancelFunc()
	})
	//这里ctx感知到cancel则会关闭watcher
	watchRespChan = watcher.Watch(ctx, "/cron/jobs/job7", clientv3.WithRev(watchStartRevision))

	// 处理kv变化事件
	for watchResp = range watchRespChan {
		for _, event = range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为:", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了", "Revision:", event.Kv.ModRevision)
			}
		}
	}
}
