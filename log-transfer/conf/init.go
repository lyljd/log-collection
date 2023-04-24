package conf

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"strings"
)

var Cfg = new(Config)

type Config struct {
	Kafka         Kafka
	Etcd          Etcd
	ElasticSearch ElasticSearch
}

type Kafka struct {
	Addr          string `envconfig:"KAFKA_ADDR"`
	ConsumerGroup string `envconfig:"KAFKA_CONSUMER_GROUP"`
	DlqTopic      string `envconfig:"DLQ_TOPIC" default:"dlq"`
}

type Etcd struct {
	Addr                string `envconfig:"ETCD_ADDR"`
	TimeOut             int    `envconfig:"ETCD_TIMEOUT" default:"3"`
	LogConfigurationKey string `envconfig:"ETCD_LOG_CONFIGURATION_KEY"`
}

type ElasticSearch struct {
	Addr string `envconfig:"ELASTICSEARCH_ADDR"`
}

func Init() {
	if err := godotenv.Load("./log-transfer/.env"); err != nil {
		panic("读取环境变量失败！" + err.Error())
	}

	if err := envconfig.Process("", Cfg); err != nil {
		panic("环境变量绑定Cfg失败！" + err.Error())
	}

	if !strings.HasPrefix(Cfg.ElasticSearch.Addr, "http://") && !strings.HasPrefix(Cfg.ElasticSearch.Addr, "https://") {
		panic("ELASTICSEARCH_ADDR must start with http:// or https://")
	}
}
