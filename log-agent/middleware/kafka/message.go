package kafka

import (
	"github.com/Shopify/sarama"
)

func SendMessage(topic, data string) error {
	_, _, err := Client.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(data),
	})
	return err
}
