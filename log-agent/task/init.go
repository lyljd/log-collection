package task

import (
	"context"
	"encoding/json"
	"github.com/hpcloud/tail"
	"log-collection/log-agent/conf"
	"log-collection/log-agent/logx"
	"log-collection/log-agent/middleware/etcd"
	"time"
)

var (
	Tasks    []*Task
	TasksMap map[string]*Task
)

type Task struct {
	Topic string `json:"topic"`
	Path  string `json:"path"`
	Line  <-chan *tail.Line
	Over  chan struct{}
}

func Init() {
	logKey := conf.Cfg.Etcd.LogConfigurationKey
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	getResp, err := etcd.Client.Get(ctx, logKey)
	if err != nil {
		panic("etcd中" + logKey + "读取失败！" + err.Error())
	}
	getRespKvs := getResp.Kvs
	if getRespKvs == nil {
		logx.Log.Println("etcd中" + logKey + "未配置")
	} else {
		load(getRespKvs[0].Value)
	}
}

func load(data []byte) {
	if err := json.Unmarshal(data, &Tasks); err != nil {
		logx.Log.Println("etcd中" + conf.Cfg.Etcd.LogConfigurationKey + "配置有误！" + err.Error())
		return
	}

	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	TasksMap = make(map[string]*Task)
	for k, v := range Tasks {
		t, err := tail.TailFile(v.Path, config)
		if err != nil {
			panic("打开日志文件失败！" + err.Error())
		}
		Tasks[k].Line = t.Lines
		Tasks[k].Over = make(chan struct{})
		TasksMap[v.Topic+":"+v.Path] = v
	}
}
