/*
@Date: 2022/4/27 13:14
@Author: max.liu
@File : consume
@Desc:
*/

package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/maxliu9403/common/kafka"
	"github.com/maxliu9403/common/logger"
)

const (
	Topic = "demo"
	group = "template-group"
)

type demoData struct {
	User string `json:"user"`
}

// RunConsume 在后台运行消费者的程序
func RunConsume(consumer *kafka.ConsumerClient) (err error) {
	hand := kafka.NewConsumerGroup(taskHandler)
	err = consumer.RunConsumer(group, []string{Topic}, hand)
	if err != nil {
		logger.Errorf("failed to consume: %s", err.Error())
		return err
	}

	if consumer.IsRunning() {
		logger.Info("kafka consumer is ready")
	} else {
		return fmt.Errorf("consumer is not running")
	}

	return err
}

func taskHandler(message *sarama.ConsumerMessage) {
	switch message.Topic {
	case Topic:
		var tmp demoData
		err := json.Unmarshal(message.Value, &tmp)
		if err != nil {
			logger.Errorf("unmarshal message of offset %d with partition %d failed: %s", message.Offset, message.Partition, err.Error())
		} else {
			err = consumeTask(tmp)
			if err != nil {
				logger.Errorf("consume data failed, message offset is %d with partition %d: %s", message.Offset, message.Partition, err.Error())
			}
		}
	default:
		logger.Errorf("unknown topic: %s", message.Topic)
	}
}

func consumeTask(data demoData) (err error) {
	logger.Debugf("process message here: %s", data.User)
	if data.User == "" {
		return fmt.Errorf("no user")
	}

	return
}
