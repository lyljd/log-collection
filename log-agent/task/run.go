package task

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"log-collection/log-agent/conf"
	"log-collection/log-agent/logx"
	"log-collection/log-agent/middleware/etcd"
	"log-collection/log-agent/middleware/kafka"
)

func Run() {
	for _, t := range Tasks {
		t.run()
	}

	watchLogConfigurationKey()
}

func (t *Task) run() {
	go func() {
		logx.Log.Println(conf.Cfg.Etcd.LogConfigurationKey + "中" + t.Topic + "(" + t.Path + ")已监听")
		for {
			select {
			case line := <-t.Line:
				if err := kafka.SendMessage(t.Topic, line.Text); err != nil {
					logx.Log.Println("向Kafka发送数据失败！" + err.Error())
				}
			case <-t.Over:
				goto outer
			}
		}
	outer:
		logx.Log.Println(conf.Cfg.Etcd.LogConfigurationKey + "中" + t.Topic + "(" + t.Path + ")已取消监听")
	}()
}

func (t *Task) stop() {
	t.Over <- struct{}{}
}

func watchLogConfigurationKey() {
	go func() {
		watchChan := etcd.Client.Watch(context.Background(), conf.Cfg.Etcd.LogConfigurationKey)
		for {
			select {
			case watchResp := <-watchChan:
				evt := watchResp.Events[0]
				switch evt.Type {
				case mvccpb.PUT:
					logx.Log.Println("etcd中" + conf.Cfg.Etcd.LogConfigurationKey + "已更新")
					updateTask(evt.Kv.Value)
				case mvccpb.DELETE:
					logx.Log.Println("etcd中" + conf.Cfg.Etcd.LogConfigurationKey + "已删除")
					for _, t := range Tasks {
						t.stop()
					}
				default:
					logx.Log.Println("未识别的etcd事件类型：" + evt.Type.String())
				}
			}
		}
	}()
}

func updateTask(data []byte) {
	for _, t := range TasksMap {
		t.stop()
	}

	load(data)

	for _, t := range TasksMap {
		t.run()
	}
}
