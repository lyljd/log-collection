package task

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"log-agent/conf"
	"log-agent/logx"
	"log-agent/middleware/etcd"
	"log-agent/middleware/kafka"
)

func Run() {
	for _, t := range Tasks {
		t.run()
	}

	watchLogConfigurationKey()
}

func (t *Task) run() {
	go func() {
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
	}()
}

func (t *Task) stop() {
	t.Over <- struct{}{}
}

func watchLogConfigurationKey() {
	go func() {
		watchChan := etcd.Client.Watch(context.Background(), conf.Cfg.Etcd.LogConfigurationKey)
		for watchResp := range watchChan {
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
				Tasks = nil
				TasksMap = make(map[string]*Task)
			default:
				logx.Log.Println("未识别的etcd事件类型：" + evt.Type.String())
			}

		}
	}()
}

func updateTask(data []byte) {
	oldTasksMap := make(map[string]*Task)
	for k, v := range TasksMap {
		oldTasksMap[k] = v
	}

	load(data)

	for k, v := range TasksMap {
		if _, ok := oldTasksMap[k]; !ok {
			v.run()
		}
	}

	for k, v := range oldTasksMap {
		if _, ok := TasksMap[k]; !ok {
			v.stop()
		}
	}
}
