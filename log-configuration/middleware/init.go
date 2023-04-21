package middleware

import "log-collection/log-configuration/middleware/etcd"

func Init() {
	etcd.Init()
}
