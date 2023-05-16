package key

import (
	"log-configuration/conf"
	"log-configuration/database/redis"
	"log-configuration/serializer"
)

type DeleteKeyService struct {
}

func (s *DeleteKeyService) DeleteKey(key string) serializer.Response {
	_, err := redis.Client.SRem(conf.Cfg.Redis.RegisterKey, key).Result()
	if err != nil {
		return serializer.SerErr("", err)
	}
	return serializer.Response{}
}
