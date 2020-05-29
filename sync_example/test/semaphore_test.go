package test_test

import (
	. ".."
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)


func TestSemaphore(t *testing.T) {
	wg:=sync.WaitGroup{}
	N:=10
	wg.Add(N)
	s:=NewSemaphore(3)
	for i:=0;i<N;i++{
		go func(i int) {
			defer wg.Done()
			s.Acquire()
			fmt.Printf("协程%d,进入，当前已有%d个并发\n",i,s.Permits())
			//设置随机种子
			rand.Seed(time.Now().UnixNano())
			time.Sleep(time.Duration(rand.Intn(10))*time.Second)	//模拟耗时操作
			fmt.Printf("协程%d,即将离开\n",i)
			s.Release()
		}(i)
	}
	wg.Wait()
}