package main

import (
	"log-collection/log-agent/conf"
	"log-collection/log-agent/database"
	"log-collection/log-agent/database/redis"
	"log-collection/log-agent/logx"
	"log-collection/log-agent/middleware"
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
		logx.Log.Println("log-agent已停止并已取消所有监听")
	}
}
