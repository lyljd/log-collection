package kafka

import (
	"github.com/Shopify/sarama"
	"log-collection/log-transfer/conf"
	"strings"
)

var Client sarama.ConsumerGroup

func Init() {
	addr := strings.Split(conf.Cfg.Kafka.Addr, ",")

	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true

	cli, err := sarama.NewConsumerGroup(addr, conf.Cfg.Kafka.ConsumerGroup, cfg)
	if err != nil {
		panic("新建kafka消费者组失败！" + err.Error())
	}

	Client = cli
}
