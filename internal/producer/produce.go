/*
@Date: 2022/4/27 13:09
@Author: max.liu
@File : produce
@Desc:
*/

package producer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/maxliu9403/common/kafka"
	"github.com/maxliu9403/common/logger"
)

const (
	TaskPoolTopic = "demo"
)

var kafkaProducer *kafka.AsyncProducer

// SendMessage 用于生产消息，可变参数 keys 表示让同一个 key 的消息发送到同一个 partition，
// 比如同一个机器的消息有先后顺序，这时 keys 设置成这台机器的 IP，保证了消费时的先后顺序。
// 需要注意，keys 长度仅允许为 1
func SendMessage(msg interface{}, keys ...string) error {
	if kafkaProducer == nil {
		return fmt.Errorf("kakfa producer is not initialized yet")
	}

	producer := *kafkaProducer
	js, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	logger.Debugf("send message is %v", string(js))
	return producer.Produce(TaskPoolTopic, js, keys...)
}

// NewProducer 从 kafka client 新建一个消费者
func NewProducer(conf kafka.Config) {
	var err error

	if kafkaProducer != nil {
		return
	}

	cli := kafka.Default()
	if cli == nil {
		cli, err = conf.BuildKafka(context.TODO())
		if err != nil {
			logger.Fatal(err.Error())
			return
		}
	}

	p, err := cli.NewAsyncProducerClient()
	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	p.RunAsyncProducer()
	go func() {
		for {
			e := <-p.ProducerErrors()
			logger.Errorf("message %v produce failed", e.Msg)
		}
	}()

	kafkaProducer = &p
}
