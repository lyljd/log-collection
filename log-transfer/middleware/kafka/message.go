package kafka

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"log-transfer/logx"
	"log-transfer/middleware/elasticsearch"
	"log-transfer/tool"
	"time"
)

type message struct {
	Content   string `json:"content"`
	Date      string `json:"date"`
	Timestamp int64  `json:"timestamp"`
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

			msg := m.Value
			if !tool.IsJSON(string(m.Value)) {
				msg, _ = json.Marshal(&message{
					Content:   string(m.Value),
					Date:      time.Now().Format(time.DateTime),
					Timestamp: time.Now().Unix(),
				})
			}

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
