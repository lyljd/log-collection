package middleware

import (
	"log-collection/log-transfer/middleware/elasticsearch"
	"log-collection/log-transfer/middleware/etcd"
	"log-collection/log-transfer/middleware/kafka"
)

func Init() {
	kafka.Init()
	etcd.Init()
	elasticsearch.Init()
}
