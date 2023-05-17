package main

import (
	"log-transfer/conf"
	"log-transfer/logx"
	"log-transfer/middleware"
	"log-transfer/middleware/etcd"
	"log-transfer/middleware/kafka"
	"log-transfer/task"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	conf.Init()
	logx.Init()
	middleware.Init()
	task.Init()
}

func main() {
	task.Run()

	logx.Log.Println("log-transfer已启动")

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
	select {
	case <-interrupt:
		releaseConn()
		logx.Log.Println("log-transfer已停止")
	}
}

func releaseConn() {
	_ = kafka.Client.Close()
	_ = etcd.Client.Close()
}
