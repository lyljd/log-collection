package kafka

import (
	"github.com/Shopify/sarama"
	"log-collection/log-transfer/conf"
	"strings"
)

var (
	Client         sarama.ConsumerGroup
	client         sarama.Client
	producerClient sarama.SyncProducer
)

func Init() {
	addr := strings.Split(conf.Cfg.Kafka.Addr, ",")

	// 新建kafka消费者组
	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true
	consumerCli, err := sarama.NewConsumerGroup(addr, conf.Cfg.Kafka.ConsumerGroup, cfg)
	if err != nil {
		panic("新建kafka消费者组失败！" + err.Error())
	}
	Client = consumerCli

	// 新建kafka客户端
	cli, err := sarama.NewClient(addr, sarama.NewConfig())
	if err != nil {
		panic("新建kafka客户端失败！" + err.Error())
	}
	client = cli

	// 新建kafka生产者
	cfg = sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	cfg.Producer.Return.Successes = true
	producerCli, err := sarama.NewSyncProducer(addr, cfg)
	if err != nil {
		panic("新建kafka生产者失败！" + err.Error())
	}
	producerClient = producerCli
}

func ReleaseConn() {
	_ = Client.Close()
	_ = client.Close()
	_ = producerClient.Close()
}
