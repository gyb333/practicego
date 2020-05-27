package sync_example_test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)
/*
用通道实现信号量，控制并发个数,获取许可(Acquire())、指定时间内获取许可(TryAcquireOnTime)、释放许可(Release())

信号量Semaphore:用来控制可同时并发的线程数，其内部维护了一组虚拟许可，通过构造器指定许可的数量，
每次线程执行操作时先通过acquire方法获得许可，执行完毕再通过release方法释放许可。如果无可用许可，那么acquire方法将一直阻塞，直到其它线程释放许可。

线程池:用来控制实际工作的线程数量，通过线程复用的方式来减小内存开销。线程池可同时工作的线程数量是一定的，
超过该数量的线程需进入线程队列等待，直到有可用的工作线程来执行任务。
 */

type Semaphore struct {
	length int 		//许可数量
	channel chan struct{}
}

func NewSemaphore(length int)  *Semaphore{
	return &Semaphore{
		length: length,
		channel: make(chan struct{},length),
	}
}

func (s *Semaphore) Acquire()  {
	s.channel<- struct{}{}
}

func (s *Semaphore) Release()  {
	<-s.channel
}

func (s *Semaphore) TryAcquire() bool {
	select {
	case s.channel<-struct{}{}:
		return true
	default:
		return false
	}
}
/* 尝试指定时间内获取许可 */
func (s *Semaphore) TryAcquireOnTime(timeout time.Duration) bool {
	select {
	case s.channel<-struct{}{}:
		return true
	case <-time.After(timeout):
		return false
	}
}

/* 当前可用的许可数 */
func (s *Semaphore) AvailablePermits() int {
	return s.length - len(s.channel)
}

func (s *Semaphore) Permits() int {
	return len(s.channel)
}



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