package kafka

import (
	"github.com/Shopify/sarama"
	"log-agent/conf"
	"strings"
)

var Client sarama.SyncProducer

func Init() {
	addr := strings.Split(conf.Cfg.Kafka.Addr, ",")

	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	cfg.Producer.Return.Successes = true

	cli, err := sarama.NewSyncProducer(addr, cfg)
	if err != nil {
		panic("新建kafka生产者失败！" + err.Error())
	}

	Client = cli
}
