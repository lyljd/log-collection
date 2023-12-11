package elasticsearch

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"io"
	"log-transfer/conf"
	"strings"
)

type fieldBody struct {
	Type     string `json:"type"`
	Analyzer string `json:"analyzer"`
}

func checkIndexExist(name string) (bool, error) {
	req := esapi.IndicesExistsRequest{
		Index: []string{name},
	}

	res, err := req.Do(context.Background(), Client)
	if err != nil {
		return false, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	if res.IsError() {
		return false, nil
	}

	return true, nil
}

func createIndex(name string, strTypProperties []string) error {
	if exist, err := checkIndexExist(name); err != nil || exist {
		return err
	}

	fb := fieldBody{"text", conf.Cfg.ElasticSearch.Analyzer}

	properties := make(map[string]fieldBody)
	for _, p := range strTypProperties {
		properties[p] = fb
	}
	mappings := map[string]map[string]fieldBody{"properties": properties}
	reqBody := map[string]map[string]map[string]fieldBody{"mappings": mappings}

	reqBodyJson, _ := json.Marshal(reqBody)

	req := esapi.IndicesCreateRequest{
		Index: name,
		Body:  strings.NewReader(string(reqBodyJson)),
	}

	if _, err := req.Do(context.Background(), Client); err != nil {
		return err
	}

	return nil
}

func getStrTypProperties(jsonVal []byte) (strTypProperties []string) {
	var mapping map[string]any
	_ = json.Unmarshal(jsonVal, &mapping)

	for k, v := range mapping {
		if _, ok := v.(string); ok {
			strTypProperties = append(strTypProperties, k)
		}
	}

	return
}
