package redis

import (
	"github.com/go-redis/redis"
	"log-configuration/conf"
)

var Client *redis.Client

func Init() {
	cli := redis.NewClient(&redis.Options{
		Addr:     conf.Cfg.Redis.Addr,
		Password: conf.Cfg.Redis.Password,
	})

	_, err := cli.Ping().Result()
	if err != nil {
		panic("新建redis客户端失败！" + err.Error())
	}

	Client = cli
}
