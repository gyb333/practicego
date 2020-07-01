package test

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"log"
	"os"
	"strconv"
	"testing"
)

func TestOrderedProducer(t *testing.T)  {
	OrderedProducer()
}

func OrderedProducer ()  {
	p, _ := rocketmq.NewProducer(
		producer.WithGroupName("testGroup"),
		//producer.WithNsResovler(primitive.NewPassthroughResolver([]string{"192.168.56.100:9876"})),
		producer.WithNameServer([]string{"192.168.56.100:9876"}),
		producer.WithRetry(2),

	)
	err := p.Start()
	if err != nil {
		log.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}
	tags := []string{"TagA", "TagB", "TagC", "TagD", "TagE"};
	for i := 0; i < 100; i++  {
		 //orderId := i % 10;
		msg := &primitive.Message{
			Topic: "study" /* Topic */,
			Body:  []byte("Hello RocketMQ Go Client! " + strconv.Itoa(i)),
		}
		msg=msg.WithTag(tags[i %len(tags)]).WithKeys([]string{"KEY"+  strconv.Itoa(i)})
		res, err := p.SendSync(context.Background(), msg)
		if err != nil {
			log.Printf("send message error: %s\n", err)
		} else {
			log.Printf("send message success: result=%s\n", res.String())
		}
	}
}
