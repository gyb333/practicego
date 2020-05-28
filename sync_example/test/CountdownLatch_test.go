package test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestCountdownLatch(t *testing.T){
	N:=10
	answer :=sync.WaitGroup{}
	answer.Add(N)
	var order=sync.WaitGroup{}
	order.Add(1)
	for i:=0;i<N;i++{
		go func(i int) {
			fmt.Printf("协程%d, 正准备接受命令\n",i)
			order.Wait();
			fmt.Printf("协程%d, 已接受命令\n",i)
			time.Sleep(time.Duration(rand.Intn(10000)))
			fmt.Printf("协程%d, 回应命令处理结果\n",i)
			answer.Done();
		}(i+1)
	}
	time.Sleep(time.Duration(rand.Intn(10000)))
	fmt.Println("主协程, 即将发布命令");
	order.Done();
	fmt.Println("主协程, 已发送命令，正在等待结果");
	answer.Wait()
	fmt.Println("主协程, 已收到所有响应结果");

}
