package test

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func NewProducer(retries int,hosts ...string) (rocketmq.Producer, error) {
	return rocketmq.NewProducer(
		//producer.WithGroupName("testGroup"),
		//producer.WithNsResovler(primitive.NewPassthroughResolver([]string{"192.168.56.100:9876"})),
		producer.WithNameServer(hosts),
		producer.WithRetry(retries),
		//producer.WithQueueSelector(producer.NewManualQueueSelector()),
	)
}
