package test

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestProducer(t *testing.T)  {
	simpleProducer()
}

func simpleProducer()  {
	p, _ := rocketmq.NewProducer(
		producer.WithGroupName("testGroup"),
		//producer.WithNsResovler(primitive.NewPassthroughResolver([]string{"192.168.56.100:9876"})),
		producer.WithNameServer([]string{"192.168.56.100:9876"}),
		producer.WithRetry(2),
	)
	err := p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}

	for i := 0; i < 1; i++ {
		msg := &primitive.Message{
			Topic: "study" /* Topic */,
			Body:  []byte("Hello RocketMQ Go Client! " + strconv.Itoa(i)),
		}
		msg=msg.WithTag("Tag").WithShardingKey("ShardingKey").
			WithKeys([]string{"K1","K2"})
		msg.WithProperty("key","value")
		res, err := p.SendSync(context.Background(), msg)
		if err != nil {
			fmt.Printf("send message error: %s\n", err)
		} else {
			fmt.Printf("send message success: result=%s\n", res.String())
		}
	}
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}
}

func TestConsumer(t *testing.T)  {
	simpleConsumer()
}

func simpleConsumer()  {
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName("testGroup"),
		//consumer.WithNsResovler(primitive.NewPassthroughResolver([]string{"192.168.56.100:9876"})),
		consumer.WithNameServer([]string{"192.168.56.100:9876"}),
	)
	err := c.Subscribe("test", consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			fmt.Printf("subscribe callback: %v \n", msgs[i])
		}

		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	// Note: start after subscribe
	err = c.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	time.Sleep(time.Hour)
	err = c.Shutdown()
	if err != nil {
		fmt.Printf("shutdown Consumer error: %s", err.Error())
	}
}
