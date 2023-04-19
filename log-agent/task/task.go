package task

import (
	"encoding/json"
	"github.com/hpcloud/tail"
	"log-collection/log-agent/conf"
	"log-collection/log-agent/logx"
	"log-collection/log-agent/middleware/kafka"
)

var (
	Tasks    []Task
	TasksMap = make(map[string]Task)
)

type Task struct {
	Topic string `json:"topic"`
	Path  string `json:"path"`
	Line  <-chan *tail.Line
	Over  chan struct{}
}

func (t Task) run() {
	go func() {
		logx.Log.Println(conf.Cfg.Etcd.LogConfigurationKey + "中" + t.Topic + "(" + t.Path + ")已监听")
		for {
			select {
			case line := <-t.Line:
				kafka.SendMessage(t.Topic, line.Text)
			case <-t.Over:
				goto outer
			}
		}
	outer:
		logx.Log.Println(conf.Cfg.Etcd.LogConfigurationKey + "中" + t.Topic + "(" + t.Path + ")已取消监听")
	}()
}

func Put(data []byte) {
	// Load task
	if err := json.Unmarshal(data, &Tasks); err != nil {
		logx.Log.Println("etcd中" + conf.Cfg.Etcd.LogConfigurationKey + "配置有误！" + err.Error())
		return
	}
	loadLine()

	// Run new task
	newTasksMap := make(map[string]Task)
	for _, t := range Tasks {
		mapKey := t.Topic + ":" + t.Path
		if _, ok := TasksMap[mapKey]; !ok {
			t.run()
		}
		newTasksMap[mapKey] = t
	}

	// Delete non-existent task
	for mapKey, t := range TasksMap {
		if _, ok := newTasksMap[mapKey]; !ok {
			t.Over <- struct{}{}
		}
	}
}

func DeleteAll() {
	for _, t := range Tasks {
		t.Over <- struct{}{}
	}
}

func loadLine() {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	for k, v := range Tasks {
		t, err := tail.TailFile(v.Path, config)
		if err != nil {
			panic("打开日志文件失败！" + err.Error())
		}
		Tasks[k].Line = t.Lines
		Tasks[k].Over = make(chan struct{})
	}
}
