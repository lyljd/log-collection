package serializer

import "encoding/json"

type Configuration struct {
	Topic string `json:"topic"`
	Path  string `json:"path"`
}

func BuildConfigurationResponse(data []byte) Response {
	cs := make([]*Configuration, 0)
	if data != nil {
		if err := json.Unmarshal(data, &cs); err != nil {
			return SerErr("etcd中key配置有误！", err)
		}
	}
	return Response{
		Data: cs,
	}
}
