package pool_test

import (
	"./queue"
	. "./simple"
	"fmt"
	"runtime"
	"testing"
	"time"
)

type Score struct {
	Num int
}
func (s *Score) Do() {
	//fmt.Println("num:", s.Num)
	time.Sleep(1 * 1 * time.Second)
}


func TestSimpleWorkerPool(t *testing.T) {
	num := 100 * 100 *50
	// debug.SetMaxThreads(num + 1000) //设置最⼤线程数
	p := NewWorkerPool(num)
	p.Run()
	datanum := 100 * 100 * 100*10
	go func() {
		for i := 1; i <= datanum; i++ {
			sc := &Score{Num: i}
			p.JobQueue <- sc
		}
	}()
	for {
		fmt.Println("runtime.NumGoroutine() :", runtime.NumGoroutine())
		time.Sleep(2 * time.Second)
	}
}
func TestWorkerPool(t *testing.T) {
	num := 100 * 100 *50
	p := queue.NewWorkerPool(num)
	p.Run()
	datanum := 100 * 100 * 100*10
	go func() {
		for i := 1; i <= datanum; i++ {
			sc := &Score{Num: i}
			p.Scheduler.Submit(sc)
		}
	}()
	for {
		fmt.Println("runtime.NumGoroutine() :", runtime.NumGoroutine())
		time.Sleep(2 * time.Second)
	}
}