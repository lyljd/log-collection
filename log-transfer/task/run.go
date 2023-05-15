package task

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"log-transfer/conf"
	"log-transfer/logx"
	"log-transfer/middleware/etcd"
	"log-transfer/middleware/kafka"
)

func Run() {
	kafka.ReceiveMessage(Topics)

	watchLogConfigurationKey()
}

func watchLogConfigurationKey() {
	go func() {
		watchChan := etcd.Client.Watch(context.Background(), conf.Cfg.Etcd.LogConfigurationKey)
		for watchResp := range watchChan {
			if len(watchResp.Events) == 0 {
				continue
			}
			evt := watchResp.Events[0]
			switch evt.Type {
			case mvccpb.PUT:
				logx.Log.Println("etcd中" + conf.Cfg.Etcd.LogConfigurationKey + "已更新")
				load(evt.Kv.Value)
				kafka.ReceiveMessage(Topics)
			case mvccpb.DELETE:
				logx.Log.Println("etcd中" + conf.Cfg.Etcd.LogConfigurationKey + "已删除")
				kafka.DestroyConsumerGroups()
				Topics = []string{}
			default:
				logx.Log.Println("未识别的etcd事件类型：" + evt.Type.String())
			}

		}
	}()
}
