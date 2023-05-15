package etcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"log-transfer/conf"
	"strings"
	"time"
)

var Client *clientv3.Client

func Init() {
	addr := strings.Split(conf.Cfg.Etcd.Addr, ",")

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   addr,
		DialTimeout: time.Duration(conf.Cfg.Etcd.TimeOut) * time.Second,
	})
	if err != nil {
		panic("新建etcd客户端失败！" + err.Error())
	}

	Client = cli
}
