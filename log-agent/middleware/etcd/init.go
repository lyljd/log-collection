package etcd

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log-collection/log-agent/conf"
	"log-collection/log-agent/logx"
	"log-collection/log-agent/task"
	"strings"
	"time"
)

var client *clientv3.Client

func Init() {
	// Init etcd client
	addr := strings.Split(conf.Cfg.Etcd.Addr, ",")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   addr,
		DialTimeout: time.Duration(conf.Cfg.Etcd.TimeOut) * time.Second,
	})
	if err != nil {
		panic("新建etcd客户端失败！" + err.Error())
	}
	client = cli

	// Load configuration
	logKey := conf.Cfg.Etcd.LogConfigurationKey
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	getResp, err := client.Get(ctx, logKey)
	if err != nil {
		panic("etcd中" + logKey + "读取失败！" + err.Error())
	}
	getRespKvs := getResp.Kvs
	if getRespKvs == nil {
		logx.Log.Println("etcd中" + logKey + "未配置")
	} else {
		task.Put(getRespKvs[0].Value)
	}

	// Watch logKey
	go func() {
		watchChan := client.Watch(context.Background(), logKey)
		for {
			select {
			case watchResp := <-watchChan:
				evt := watchResp.Events[0]
				switch evt.Type {
				case mvccpb.PUT:
					logx.Log.Println("etcd中" + conf.Cfg.Etcd.LogConfigurationKey + "已更新")
					task.Put(evt.Kv.Value)
				case mvccpb.DELETE:
					logx.Log.Println("etcd中" + conf.Cfg.Etcd.LogConfigurationKey + "已删除")
					task.DeleteAll()
				default:
					logx.Log.Println("未识别的etcd事件类型：" + evt.Type.String())
				}
			}
		}
	}()
}
