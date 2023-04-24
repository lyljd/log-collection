package main

import (
	"log-collection/log-transfer/conf"
	"log-collection/log-transfer/logx"
	"log-collection/log-transfer/middleware"
	"log-collection/log-transfer/middleware/etcd"
	"log-collection/log-transfer/middleware/kafka"
	"os"
	"os/signal"
)

func init() {
	conf.Init()
	logx.Init()
	middleware.Init()
}
func main() {
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
	_ = kafka.Client.Close()
	_ = etcd.Client.Close()
}
