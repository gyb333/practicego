package ratelimit

import (
	"AdvancedResearch/ratelimit/leakybucket"
	"AdvancedResearch/ratelimit/simpleratelimit"
	"context"
	"fmt"
	"go.uber.org/ratelimit"
	"golang.org/x/time/rate"
	"log"
	"testing"
	"time"
)

func TestRatelimit(t *testing.T)  {
	// rate limit: simple
	rl := simpleratelimit.New(10, time.Second)

	for i := 0; i < 100; i++ {
		log.Printf("limit result: %v\n", rl.Limit())
	}
	log.Printf("limit result: %v\n", rl.Limit())

}
func TestLeakybucket(t *testing.T) {
	// rate limit: leaky-bucket
	lb := leakybucket.New()
	b, err := lb.Create("leaky_bucket", 10, time.Second)
	if err != nil {
		log.Println(err)
	}
	log.Printf("bucket capacity:%v", b.Capacity())
}
func TestUber(t *testing.T)  {
	// rate limit: token-bucket
	rl := ratelimit.New(100) // per second
	prev := time.Now()
	for i := 0; i < 10; i++ {
		now := rl.Take()
		if i > 0 {
			fmt.Println(i, now.Sub(prev))
		}
		prev = now
	}
}
func TestRate(t *testing.T)  {
	l := rate.NewLimiter(2, 5)
	ctx := context.Background()
	start := time.Now()
	// 要处理二十个事件
	for i := 0; i < 20; i++ {
		l.Wait(ctx)
		// dosomething
	}
	fmt.Println(time.Since(start)) // output: 7.501262697s （初始桶内5个和每秒2个token）
}
