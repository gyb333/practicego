package basic

import (
	"runtime"
	"time"
	"fmt"
)

func ChanMain()  {
	//baseChan()
	//chanDemo()
	//bufferedChannel()
	channelClose()
}


func baseChan(){
	c := make(chan struct{})
	ci := make(chan int, 100)
	go func(i chan struct{}, j chan int) {
		for i := 0; i < 10; i++ {
			j <- i
		}
		close(j)
		i <- struct{}{}
	}(c, ci)
	println("NumGoroutine=", runtime.NumGoroutine())
	<-c
	println("NumGoroutine=", runtime.NumGoroutine())
	for v := range ci {
		print(v)
	}
}

type Worker struct {
	id int
	i chan int
}

func createWorker(id int) *Worker  {
	c :=make(chan int)
	return &Worker{id:id,i:c}
}

func doWork(w *Worker)  {
	for n:=range w.i{
		fmt.Printf("Worker %d received %c\n",w.id,n)
	}
}

func chanDemo()  {
	var channels [10]*Worker
	for i=0;i< len(channels);i++{
		w:=createWorker(i)
		channels[i]=w
		go doWork(w)

	}
	for i=0;i< len(channels);i++{
		channels[i].i<-'a'+i
	}
	for i=0;i< len(channels);i++{
		channels[i].i<-'A'+i
	}

	time.Sleep(time.Microsecond)
}

func bufferedChannel(){
	c :=make(chan int ,3)
	w :=&Worker{id:1,i:c}
	go doWork(w)
	for i=0;i< 10;i++{
		w.i<-'a'+i
	}
	time.Sleep(time.Microsecond)
}

func channelClose(){
	c :=make(chan int ,3)
	w :=&Worker{id:1,i:c}
	go doWork(w)
	for i=0;i< 10;i++{
		w.i<-'a'+i
	}
	close(c)
	time.Sleep(time.Microsecond)
}
