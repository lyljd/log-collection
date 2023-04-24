package elasticsearch

import (
	"bytes"
	"context"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func SendMessage(index string, body []byte) error {
	req := esapi.IndexRequest{
		Index: index,
		Body:  bytes.NewReader(body),
	}
	if _, err := req.Do(context.Background(), Client); err != nil {
		return err
	}
	return nil
}
