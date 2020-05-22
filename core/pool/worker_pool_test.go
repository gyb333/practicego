package pool_test

import (
	"./queue"
	. "./simple"
	"sync"
	"testing"
)

type Score struct {
	Num int
	*sync.WaitGroup
}
func (s *Score) Do() {
	//fmt.Println("num:", s.Num)
	//time.Sleep(10* time.Microsecond)
	defer s.Done()
}

var num = 100*100
var datanum = 100 * 100*100*10

func simpleWorkerPool()  {
	p := NewWorkerPool(num)
	p.Run()
	wg:=sync.WaitGroup{}
	wg.Add(datanum)
	go func() {
		for i := 1; i <= datanum; i++ {
			sc := &Score{Num: i,WaitGroup: &wg}
			p.JobQueue <- sc
		}
	}()
	wg.Wait()
}

func TestSimpleWorkerPool(t *testing.T) {
	simpleWorkerPool()
}

func BenchmarkSimpleWorkerPool(b *testing.B) {
	simpleWorkerPool()
}

func queueWorkerPool()  {
	p := queue.NewWorkerPool(num)
	p.Run()
	wg:=sync.WaitGroup{}
	wg.Add(datanum)
	go func() {
		for i := 1; i <= datanum; i++ {
			sc := &Score{Num: i,WaitGroup:&wg}
			p.Scheduler.Submit(sc)
		}
	}()
	wg.Wait()
}

func TestWorkerPool(t *testing.T) {
	queueWorkerPool()
}

func BenchmarkWorkerPool(b *testing.B) {
	queueWorkerPool()
}