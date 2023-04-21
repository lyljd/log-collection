package serializer

type Key struct {
	Key string `json:"key"`
}

func buildKey(key string) Key {
	return Key{
		Key: key,
	}
}

func BuildKeysResponse(keys []string) Response {
	Keys := make([]Key, len(keys))
	for i, key := range keys {
		Keys[i] = buildKey(key)
	}
	return Response{
		Data: Keys,
	}
}
