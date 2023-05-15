package middleware

import "log-configuration/middleware/etcd"

func Init() {
	etcd.Init()
}
