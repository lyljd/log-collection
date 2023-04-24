package kafka

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"log-collection/log-transfer/conf"
	"log-collection/log-transfer/logx"
	"log-collection/log-transfer/middleware/elasticsearch"
	"runtime"
)

var consumerGroups []*consumerGroup

type consumerGroup struct {
	over chan struct{}
}

type message struct {
	Log  string `json:"log"`
	Time string `json:"time"`
}

func (g *consumerGroup) Setup(sarama.ConsumerGroupSession) error { return nil }

func (g *consumerGroup) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (g *consumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case m := <-claim.Messages():
			if m == nil {
				goto outer
			}
			msg, err := json.Marshal(&message{
				Log:  string(m.Value),
				Time: m.Timestamp.Format("2006-01-02 15:04:05"),
			})
			if err != nil {
				logx.Log.Println("将message json化失败！" + err.Error())
			} else if err = elasticsearch.SendMessage(m.Topic, msg); err != nil {
				logx.Log.Printf("向elasticsearch发送数据失败！%s；此条数据的topic为：%s，msg为：\n%s\n", err.Error(), m.Topic, string(msg))
				if err := sendMessage(conf.Cfg.Kafka.DlqTopic, string(msg)); err != nil {
					logx.Log.Println("此条数据发往死信队列失败")
				} else {
					logx.Log.Println("此条数据已发往死信队列(topic: " + conf.Cfg.Kafka.DlqTopic + ")，请在kafka日志过期前处理")
				}
			}
			session.MarkMessage(m, "")
		case <-g.over:
			goto outer
		}
	}
outer:
	return nil
}

func ReceiveMessage(topic []string) {
	DestroyConsumerGroups()

	listeningError()

	cn := getConsumerNum(topic)
	failNum := 0
	for i := 0; i < cn; i++ {
		go func() {
			err := Client.Consume(context.Background(), topic, &consumerGroup{
				over: make(chan struct{}),
			})
			if err != nil {
				failNum++
				logx.Log.Printf("新建kafka消费者失败(%d/%d)！%s\n", failNum, cn, err.Error())
			}
		}()
	}
	if failNum == cn {
		logx.Log.Println("新建kafka消费者全部失败，log-transfer已失效！")
		return
	}
	if failNum > cn/2 {
		logx.Log.Println("新建kafka消费者失败数超过一办，消费性能可能会受影响！")
	}
}

func getConsumerNum(topic []string) (consumerNum int) {
	for _, t := range topic {
		if ps, _ := client.Partitions(t); len(ps) > consumerNum {
			consumerNum = len(ps)
		}
	}
	if v := runtime.NumCPU(); consumerNum > v {
		consumerNum = v
	}
	return
}

func listeningError() {
	go func() {
		for err := range Client.Errors() {
			logx.Log.Println("kafka发生错误！" + err.Error())
		}
	}()
}

func DestroyConsumerGroups() {
	for _, g := range consumerGroups {
		g.over <- struct{}{}
	}
}

func sendMessage(topic, data string) error {
	_, _, err := producerClient.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(data),
	})
	return err
}
