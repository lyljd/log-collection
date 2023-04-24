package main

import (
	"log-collection/log-agent/conf"
	"log-collection/log-agent/database"
	"log-collection/log-agent/database/redis"
	"log-collection/log-agent/logx"
	"log-collection/log-agent/middleware"
	"log-collection/log-agent/middleware/etcd"
	"log-collection/log-agent/middleware/kafka"
	"log-collection/log-agent/task"
	"os"
	"os/signal"
)

func init() {
	conf.Init()
	logx.Init()
	database.Init()
	middleware.Init()
	task.Init()
}

func main() {
	if err := redis.RegisterService(); err != nil {
		logx.Log.Fatalln(err)
	}

	task.Run()

	logx.Log.Println("log-agent已启动")

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	select {
	case <-interrupt:
		redis.DestroyService()
		releaseConn()
		logx.Log.Println("log-agent已停止")
	}
}

func releaseConn() {
	_ = redis.Client.Close()
	_ = kafka.Client.Close()
	_ = etcd.Client.Close()
}
