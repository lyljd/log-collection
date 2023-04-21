package main

import (
	"log-collection/log-configuration/conf"
	"log-collection/log-configuration/database"
	"log-collection/log-configuration/logx"
	"log-collection/log-configuration/middleware"
	"log-collection/log-configuration/server"
	"os"
	"os/signal"
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
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	select {
	case <-interrupt:
		logx.Log.Println("log-configuration已停止")
	}
}
