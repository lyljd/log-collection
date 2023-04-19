package conf

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var Cfg = new(Config)

type Config struct {
	Kafka Kafka
	Etcd  Etcd
}

type Kafka struct {
	Addr string `envconfig:"KAFKA_ADDR"`
}

type Etcd struct {
	Addr                string `envconfig:"ETCD_ADDR"`
	TimeOut             int    `envconfig:"ETCD_TIMEOUT" default:"3"`
	LogConfigurationKey string `envconfig:"ETCD_LOG_CONFIGURATION_KEY"`
}

func Init() {
	if err := godotenv.Load("./log-agent/.env"); err != nil {
		panic("读取环境变量失败！" + err.Error())
	}

	if err := envconfig.Process("", Cfg); err != nil {
		panic("环境变量绑定Cfg失败！" + err.Error())
	}
}