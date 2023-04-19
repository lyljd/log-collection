package kafka

import (
	"github.com/Shopify/sarama"
	"log-collection/log-agent/conf"
	"strings"
)

var client sarama.SyncProducer

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

	client = cli
}