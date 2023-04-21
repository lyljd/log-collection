package conf

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var Cfg = new(Config)

type Config struct {
	Etcd  Etcd
	Redis Redis
}

type Etcd struct {
	Addr                string `envconfig:"ETCD_ADDR"`
	TimeOut             int    `envconfig:"ETCD_TIMEOUT" default:"3"`
	LogConfigurationKey string `envconfig:"ETCD_LOG_CONFIGURATION_KEY"`
}

type Redis struct {
	Addr        string `envconfig:"REDIS_ADDR"`
	Password    string `envconfig:"REDIS_PASSWORD"`
	RegisterKey string `envconfig:"REDIS_REGISTER_KEY"`
}

func Init() {
	if err := godotenv.Load("./log-configuration/.env"); err != nil {
		panic("读取环境变量失败！" + err.Error())
	}

	if err := envconfig.Process("", Cfg); err != nil {
		panic("环境变量绑定Cfg失败！" + err.Error())
	}
}
