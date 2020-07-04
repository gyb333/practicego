package example_test

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"log"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

/*
收发普通消息:是指RocketMQ 中无特性的消息，区别于有特性的定时消息、顺序消息和事务消息。
 */
const (
	Group="rocketmq_group"
	Topic="rocketmq_Topic"
	NameServer="192.168.56.100:9876"
)

func TestSimpleProducer(t *testing.T)  {
	p, _ := rocketmq.NewProducer(
		producer.WithGroupName(Group),
		//producer.WithNsResovler(primitive.NewPassthroughResolver([]string{"192.168.56.100:9876"})),
		producer.WithNameServer([]string{NameServer}),
		producer.WithRetry(2),
	)
	err := p.Start()
	if err != nil {
		log.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}

	for i := 0; i < 10; i++ {
		msg := &primitive.Message{
			Topic: Topic /* Topic */,
			Body:  []byte("Hello RocketMQ Go Client! " + strconv.Itoa(i))	/* Body */,
		}
		msg=msg.WithTag("SimpleTag")	/* Tag */
			//WithShardingKey("ShardingKey").	顺序
			//WithKeys([]string{"K1","K2"})	/* Key */
			//msg.WithProperty("key","value")
		res, err := p.SendSync(context.Background(), msg)
		if err != nil {
			log.Printf("send message error: %s\n", err)
		} else {
			log.Printf("send message success: result=%s\n", res.String())
		}
	}
	err = p.Shutdown()
	if err != nil {
		log.Printf("shutdown producer error: %s", err.Error())
	}
}

func TestSimpleAsyncProducer(t *testing.T)  {
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

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		msg := &primitive.Message{
			Topic: Topic /* Topic */,
			Body:  []byte("Hello RocketMQ Go Client! " + strconv.Itoa(i))	/* Body */,
		}
		msg=msg.WithTag("AsyncTag"/* Tag */).
				WithKeys([]string{"Async"}/* Key */)
		err := p.SendAsync(context.Background(),
			func(ctx context.Context, result *primitive.SendResult, e error) {
				if e != nil {
					log.Printf("receive message error: %s\n", err)
				} else {
					log.Printf("send message success: result=%s\n", result.String())
				}
				wg.Done()
			}, msg)

		if err != nil {
			log.Printf("send message error: %s\n", err)
		}
	}
	wg.Wait()
	err = p.Shutdown()
	if err != nil {
		log.Printf("shutdown producer error: %s", err.Error())
	}
}

func TestOneWayProducer(t *testing.T)  {
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

	for i := 0; i < 10; i++ {
		msg := &primitive.Message{
			Topic: Topic /* Topic */,
			Body:  []byte("Hello RocketMQ Go Client! " + strconv.Itoa(i))	/* Body */,
		}
		msg=msg.WithTag("OneWayTag")	/* Tag */
		err := p.SendOneWay(context.Background(), msg)
		if err != nil {
			log.Printf("send message error: %s\n", err)
		}
	}
	err = p.Shutdown()
	if err != nil {
		log.Printf("shutdown producer error: %s", err.Error())
	}
}

func TestSimpleConsumer(t *testing.T)  {
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName(Group),
		//consumer.WithNsResovler(primitive.NewPassthroughResolver([]string{"192.168.56.100:9876"})),
		consumer.WithNameServer([]string{NameServer}),
	)
	err := c.Subscribe(Topic, consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			log.Printf("subscribe callback: %v \n", msgs[i])
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
