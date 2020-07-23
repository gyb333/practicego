package pubsub

import (
	"fmt"
	"strings"
	"testing"
	"time"
)



func TestProducerConsumer(t *testing.T) {
	ch := make(chan int, 64) // 成果队列

	go Producer(3, ch) // 生成 3 的倍数的序列
	go Producer(5, ch) // 生成 5 的倍数的序列
	go Consumer(ch)    // 消费 生成的队列

	// 运行一定时间后退出
	time.Sleep(5 * time.Second)
}

func TestPubSub(t *testing.T) {
	p := NewPublisher(100*time.Millisecond, 10)
	defer p.Close()

	all := p.Subscribe()
	golang := p.SubscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "golang")
		}
		return false
	})

	p.Publish("hello,  world!")
	p.Publish("hello, golang!")

	go func() {
		for  msg := range all {
			fmt.Println("all:", msg)
		}
	} ()

	go func() {
		for  msg := range golang {
			fmt.Println("golang:", msg)
		}
	} ()

	// 运行一定时间后退出
	time.Sleep(3 * time.Second)
}