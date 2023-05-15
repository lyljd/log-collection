package configuration

import (
	"context"
	"log-configuration/middleware/etcd"
	"log-configuration/serializer"
	"time"
)

type GetConfigurationService struct {
}

func (s *GetConfigurationService) GetConfigurationByKey(key string) serializer.Response {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	getResp, err := etcd.Client.Get(ctx, key)
	if err != nil {
		return serializer.SerErr("", err)
	}
	var result []byte
	if v := getResp.Kvs; v != nil {
		result = v[0].Value
	}
	return serializer.BuildConfigurationResponse(result)
}
