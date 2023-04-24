package main

import (
	"log-collection/log-transfer/conf"
	"log-collection/log-transfer/logx"
	"log-collection/log-transfer/middleware"
	"log-collection/log-transfer/middleware/etcd"
	"log-collection/log-transfer/middleware/kafka"
	"log-collection/log-transfer/task"
	"os"
	"os/signal"
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
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	select {
	case <-interrupt:
		releaseConn()
		logx.Log.Println("log-transfer已停止")
	}
}

func releaseConn() {
	kafka.ReleaseConn()
	_ = etcd.Client.Close()
}
