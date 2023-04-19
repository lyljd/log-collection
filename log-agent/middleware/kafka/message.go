package kafka

import (
	"github.com/Shopify/sarama"
	"log-collection/log-agent/logx"
)

func SendMessage(topic, data string) {
	partition, offset, err := client.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(data),
	})
	if err != nil {
		logx.Log.Println("向kafka发送数据失败", "partition:", partition, "offset:", offset, "err:", err)
	}
}
