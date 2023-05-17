package main

import (
	"log-configuration/conf"
	"log-configuration/database"
	"log-configuration/database/redis"
	"log-configuration/logx"
	"log-configuration/middleware"
	"log-configuration/middleware/etcd"
	"log-configuration/server"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	conf.Init()
	logx.Init()
	database.Init()
	middleware.Init()
}

func main() {
	server.Run()

	logx.Log.Println("log-configuration已启动")

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
	select {
	case <-interrupt:
		releaseConn()
		logx.Log.Println("log-configuration已停止")
	}
}

func releaseConn() {
	_ = redis.Client.Close()
	_ = etcd.Client.Close()
}
