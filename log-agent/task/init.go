package task

import (
	"context"
	"encoding/json"
	"github.com/hpcloud/tail"
	"log-agent/conf"
	"log-agent/logx"
	"log-agent/middleware/etcd"
	"time"
)

var (
	Tasks    []*Task
	TasksMap = make(map[string]*Task)
	pathsMap = make(map[string]*Task)
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
		ReOpen:    false,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	TasksMap = make(map[string]*Task)
	newPathsMap := make(map[string]*Task)
	for k, v := range Tasks {
		if t, ok := pathsMap[v.Path]; ok {
			Tasks[k].Line = t.Line
		} else {
			t, err := tail.TailFile(v.Path, config)
			if err != nil {
				logx.Log.Println("打开日志文件失败！" + err.Error())
			}
			Tasks[k].Line = t.Lines
		}
		Tasks[k].Over = make(chan struct{})
		TasksMap[v.Topic+":"+v.Path] = v
		newPathsMap[v.Path] = v
	}
	pathsMap = newPathsMap
}
