package example_test

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestOrderedProducer(t *testing.T)  {
	p, _ := rocketmq.NewProducer(
		producer.WithGroupName(Group),
		producer.WithNameServer([]string{NameServer}),
		producer.WithRetry(2),
	)
	err := p.Start()
	if err != nil {
		log.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}
	tags := []string{"TagA", "TagB", "TagC", "TagD", "TagE"};
	for i := 0; i < 100; i++  {
		orderId := i % 10;
		msg := &primitive.Message{
			Topic: "OrderedTopic" /* Topic */,
			Body:  []byte("Hello RocketMQ Go Client! " + strconv.Itoa(i)),/* Body */
		}
		msg=msg.WithTag(tags[i %len(tags)]/* Tag */).
			WithKeys([]string{"KEY"+  strconv.Itoa(i)}/* Key */).
			WithShardingKey(strconv.Itoa(orderId)/*orderID*/)
		res, err := p.SendSync(context.Background(), msg)
		if err != nil {
			log.Printf("send message error: %s\n", err)
		} else {
			log.Printf("send message success: result=%s\n", res.String())
		}
	}
}


func TestOrderedConsumer(t *testing.T)  {
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName(Group),
		consumer.WithNameServer([]string{NameServer}),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
		consumer.WithConsumerOrder(true),
	)
	selector := consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: "TagA || TagC || TagD",
	}
	err := c.Subscribe("OrderedTopic", selector, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		//orderlyCtx, _ := primitive.GetOrderlyCtx(ctx)
		//log.Printf("orderly context: %v\n", orderlyCtx)
		//log.Printf("subscribe orderly callback: %v \n", msgs)
		for i := range msgs {
			fmt.Printf("subscribe callback: %v \n", msgs[i])
		}
		return consumer.ConsumeSuccess, nil

	})
	if err != nil {
		log.Println(err.Error())
	}
	// Note: start after subscribe
	err = c.Start()
	if err != nil {
		log.Println(err.Error())
		os.Exit(-1)
	}
	time.Sleep(time.Hour)
	err = c.Shutdown()
	if err != nil {
		log.Printf("shutdown Consumer error: %s", err.Error())
	}
}