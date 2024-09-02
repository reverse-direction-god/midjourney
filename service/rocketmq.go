package service

import (
	"context"
	"encoding/json"
	"fmt"
	"mj/model"
	"strconv"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/cloudflare/cfssl/log"
)

type RocketMq struct{}

//	func init() {
//		go pullSD()
//	}
func (s *RocketMq) PushMJ(req model.Queue) error {
	by, err := json.Marshal(req)
	topic := "mj"
	// 构建一个消息
	message := primitive.NewMessage(topic, by)
	_, err = MQ.SendSync(context.Background(), message)
	if err != nil {
		log.Error(err)
		return err
	} else {
		return nil
	}
}
func (s *RocketMq) PullMJ() model.Queue {
	var queue model.Queue

	ch := make(chan bool)
	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName("mjGroup"),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"*"})),
		consumer.WithPullBatchSize(1),
	)
	if err != nil {
		log.Error(err)
		return queue
	}
	if err := c.Subscribe("mj",
		consumer.MessageSelector{},
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for i := range msgs {
				json.Unmarshal(msgs[i].Message.Body, &queue)
				fmt.Println("获取到：：：", string(msgs[i].Message.Body))
				ch <- true
				break
			}
			return consumer.ConsumeSuccess, nil
		}); err != nil {
		log.Error(err)
		return queue
	}
	err = c.Start()
	if err != nil {
		log.Error(err)
		return queue
	}

	select {
	case <-ch:
		c.Shutdown()
		return queue
	}
}
func (s *RocketMq) PushSD(req model.SDPromptConfig) error {

	by, err := json.Marshal(req)
	topic := "threezto-test"
	// 构建一个消息

	message := primitive.NewMessage(topic, by)
	_, err = MQ.SendSync(context.Background(), message)
	if err != nil {
		log.Error(err)
		return err
	} else {
		return nil
	}

}

func pullSD() {

	var sd SD

	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName("resp"),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"*.*.*.*:*"})),
		consumer.WithPullBatchSize(1),
	)
	if err != nil {
		panic(err)
	}

	if err := c.Subscribe("threezto-test",
		consumer.MessageSelector{},
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for i := range msgs {

				var mod model.SDPromptConfig
				var sqlCreate model.Message
				json.Unmarshal(msgs[i].Message.Body, &mod)
				stringss, err := sd.Txt2Img(mod, mod.UserId)
				sqlCreate.ID = 0
				sqlCreate.UserId, _ = strconv.Atoi(mod.UserId)
				if err != nil {
					log.Error(err)
					sqlCreate.ResponseMessage = "err"
					sqlCreate.Url = ""
					sqlCreate.RequestMessage = mod.Prompt

				} else {
					sqlCreate.ResponseMessage = mod.Prompt
					sqlCreate.Url = stringss
					sqlCreate.RequestMessage = ""
				}
				DB.Model(&model.Message{}).Create(&sqlCreate)
				fmt.Println("收到的消息有", string(msgs[i].Message.Body))

			}
			return consumer.ConsumeSuccess, nil
		}); err != nil {
		panic(err)
	}
	err = c.Start()
	if err != nil {
		panic("启动consumer失败")
	}
	select {}
}
