package middleware

import (
	"log-agent/middleware/etcd"
	"log-agent/middleware/kafka"
)

func Init() {
	kafka.Init()
	etcd.Init()
}
