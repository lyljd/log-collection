package task

import (
	"context"
	"encoding/json"
	"log-transfer/conf"
	"log-transfer/logx"
	"log-transfer/middleware/etcd"
	"time"
)

var topics []string

type cfgItem struct {
	Topic string `json:"topic"`
}

func Init() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	logKey := conf.Cfg.Etcd.LogConfigurationKey
	getResp, err := etcd.Client.Get(ctx, logKey)
	if err != nil {
		panic("etcd中" + logKey + "读取失败！" + err.Error())
	}

	getRespKvs := getResp.Kvs
	if getRespKvs == nil {
		logx.Log.Println("etcd中" + logKey + "未配置")
		return
	}

	loadTopics(getRespKvs[0].Value)
}

func loadTopics(cfgJson []byte) {
	var cfgs []cfgItem
	if err := json.Unmarshal(cfgJson, &cfgs); err != nil {
		logx.Log.Println("etcd中" + conf.Cfg.Etcd.LogConfigurationKey + "配置有误！" + err.Error())
		return
	}

	topics = make([]string, len(cfgs))
	for k, d := range cfgs {
		topics[k] = d.Topic
	}
}
