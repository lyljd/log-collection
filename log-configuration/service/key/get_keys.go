package key

import (
	"log-collection/log-configuration/conf"
	"log-collection/log-configuration/database/redis"
	"log-collection/log-configuration/serializer"
)

type GetKeysService struct {
}

func (s *GetKeysService) GetKeys() serializer.Response {
	result, err := redis.Client.SMembers(conf.Cfg.Redis.RegisterKey).Result()
	if err != nil {
		return serializer.SerErr("", err)
	}
	return serializer.BuildKeysResponse(result)
}
