package main

import (
	"log-collection/log-agent/conf"
	"log-collection/log-agent/logx"
	"log-collection/log-agent/middleware"
	"os"
	"os/signal"
)

func init() {
	conf.Init()
	logx.Init()
	middleware.Init()
}

func main() {
	logx.Log.Println("log-agent已启动")

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	select {
	case <-interrupt:
		logx.Log.Println("log-agent已停止并已取消所有监听")
	}
}
