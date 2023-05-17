package kafka

import (
	"github.com/Shopify/sarama"
	"log-transfer/conf"
	"strings"
)

var Client sarama.ConsumerGroup

func Init() {
	addr := strings.Split(conf.Cfg.Kafka.Addr, ",")

	cli, err := sarama.NewConsumerGroup(addr, conf.Cfg.Kafka.ConsumerGroup, sarama.NewConfig())
	if err != nil {
		panic("新建kafka消费者组失败！" + err.Error())
	}

	Client = cli
}
