package task

import (
	"context"
	"encoding/json"
	"log-collection/log-transfer/conf"
	"log-collection/log-transfer/logx"
	"log-collection/log-transfer/middleware/etcd"
	"time"
)

var Topics []string

type data struct {
	Topic string `json:"topic"`
}

func Init() {
	logKey := conf.Cfg.Etcd.LogConfigurationKey
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	getResp, err := etcd.Client.Get(ctx, logKey)
	if err != nil {
		panic("etcd中" + logKey + "读取失败！" + err.Error())
	}
	getRespKvs := getResp.Kvs
	if getRespKvs == nil {
		logx.Log.Println("etcd中" + logKey + "未配置")
	} else {
		load(getRespKvs[0].Value)
	}
}

func load(value []byte) {
	var datum []*data
	if err := json.Unmarshal(value, &datum); err != nil {
		logx.Log.Println("etcd中" + conf.Cfg.Etcd.LogConfigurationKey + "配置有误！" + err.Error())
		return
	}

	Topics = make([]string, len(datum))
	for k, d := range datum {
		Topics[k] = d.Topic
	}
}
