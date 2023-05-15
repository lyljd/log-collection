package redis

import (
	"errors"
	"log-agent/conf"
)

func RegisterService() error {
	k, v := conf.Cfg.Redis.RegisterKey, conf.Cfg.Etcd.LogConfigurationKey

	ok, err := Client.SIsMember(k, v).Result()
	if err != nil {
		return err
	}
	if ok {
		return errors.New("redis注册服务失败，在" + k + "中" + v + "已存在！")
	}

	_, err = Client.SAdd(k, v).Result()
	if err != nil {
		return err
	}

	return nil
}

func DestroyService() {
	Client.SRem(conf.Cfg.Redis.RegisterKey, conf.Cfg.Etcd.LogConfigurationKey)
}
