package pool_test

import (
	"./queue"
	. "./simple"
	"fmt"
	"sync"
	"testing"
)

type Score struct {
	Num int
	sync.WaitGroup
}
func (s *Score) Do() {
	fmt.Println("num:", s.Num)
	//time.Sleep(1 * 1 * time.Second)
	defer s.Done()
}


func TestSimpleWorkerPool(t *testing.T) {
	num := 100
	// debug.SetMaxThreads(num + 1000) //设置最⼤线程数
	p := NewWorkerPool(num)
	p.Run()
	datanum := 100 * 10
	wg:=sync.WaitGroup{}
	wg.Add(datanum)
	go func() {
		for i := 1; i <= datanum; i++ {
			sc := &Score{Num: i,WaitGroup: wg}
			p.JobQueue <- sc
		}
	}()
	wg.Wait()
}
func TestWorkerPool(t *testing.T) {
	num := 100
	p := queue.NewWorkerPool(num)
	p.Run()
	datanum := 100 * 10
	wg:=sync.WaitGroup{}
	wg.Add(datanum)
	go func() {
		for i := 1; i <= datanum; i++ {
			sc := &Score{Num: i,WaitGroup:wg}
			p.Scheduler.Submit(sc)
		}
	}()
	wg.Wait()
}