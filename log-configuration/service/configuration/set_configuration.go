package configuration

import (
	"context"
	"log-configuration/middleware/etcd"
	"log-configuration/serializer"
	"time"
)

type SetConfigurationService struct {
	Data string `json:"data"`
}

func (s *SetConfigurationService) SetConfigurationByKey(key string) serializer.Response {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if s.Data == "" {
		_, err := etcd.Client.Delete(ctx, key)
		if err != nil {
			return serializer.SerErr("", err)
		}
	} else {
		_, err := etcd.Client.Put(ctx, key, s.Data)
		if err != nil {
			return serializer.SerErr("", err)
		}
	}

	return serializer.Response{}
}
