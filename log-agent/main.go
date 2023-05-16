package main

import (
	"log-agent/conf"
	"log-agent/database"
	"log-agent/database/redis"
	"log-agent/logx"
	"log-agent/middleware"
	"log-agent/middleware/etcd"
	"log-agent/middleware/kafka"
	"log-agent/task"
	"os"
	"os/signal"
	"syscall"
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
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
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
