package kafka

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"log-transfer/logx"
	"log-transfer/middleware/elasticsearch"
)

type message struct {
	Log  string `json:"log"`
	Time string `json:"time"`
}

type consumerGroup struct{}

func (g *consumerGroup) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (g *consumerGroup) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (g *consumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case m := <-claim.Messages():
			if m == nil {
				goto outer
			}
			msg, _ := json.Marshal(&message{
				Log:  string(m.Value),
				Time: m.Timestamp.Format("2006-01-02 15:04:05"),
			})
			if err := elasticsearch.SendMessage(m.Topic, msg); err != nil {
				logx.Log.Println("向es发送数据失败！" + err.Error())
				logx.Log.Println("此条数据为：" + string(msg))
			}
			session.MarkMessage(m, "")
		}
	}
outer:
	return nil
}

func StartReceiveMessage(topics []string) {
	if len(topics) == 0 {
		logx.Log.Println("开始接收消息失败！因为欲接收的消息主题列表为空。")
		return
	}
	go func() {
		// Consume()函数实际上是调用ConsumeClaim()方法
		if err := Client.Consume(context.Background(), topics, &consumerGroup{}); err != nil {
			logx.Log.Println("从Kafka接收数据失败！" + err.Error())
		}
	}()
}
