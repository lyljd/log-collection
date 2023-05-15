package middleware

import (
	"log-transfer/middleware/elasticsearch"
	"log-transfer/middleware/etcd"
	"log-transfer/middleware/kafka"
)

func Init() {
	kafka.Init()
	etcd.Init()
	elasticsearch.Init()
}
