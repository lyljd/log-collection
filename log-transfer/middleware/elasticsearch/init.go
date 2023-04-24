package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v7"
	"log-collection/log-transfer/conf"
	"strings"
)

var Client *elasticsearch.Client

func Init() {
	addr := strings.Split(conf.Cfg.ElasticSearch.Addr, ",")

	cfg := elasticsearch.Config{
		Addresses: addr,
	}

	cli, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic("新建elasticsearch客户端失败！" + err.Error())
	}

	Client = cli
}
