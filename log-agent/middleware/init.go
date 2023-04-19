package middleware

import (
	"log-collection/log-agent/middleware/etcd"
	"log-collection/log-agent/middleware/kafka"
)

func Init() {
	kafka.Init()
	etcd.Init()
}
