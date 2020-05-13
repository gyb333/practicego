package basic

import (
	"time"
	"math/rand"
	"fmt"
)

func SelectMain()  {
	var c1 ,c2 =generator(),generator()
	var worker =createWork(0)
	var values []int
	tm :=time.After(10*time.Second)
	tick :=time.Tick(time.Second)
	var activeWorker chan<-int
	var activeValue int

	for{
		if len(values)>0{
			activeWorker=worker
			activeValue=values[0]
		}
		select{
			case n:=<-c1:
				values=append(values,n)
				//fmt.Println("<-c1")
			case n:=<-c2:
				values=append(values,n)
				//fmt.Println("<-c2")
			case activeWorker <-activeValue:
				values=values[1:]
			case <-time.After(300*time.Millisecond):
				fmt.Println("timeout")
			case <-tick:
				fmt.Println("queue len=",len(values))
			case <-tm:
				fmt.Println("bye")
				return

		}
	}
}


func generator() chan int{
	out :=make(chan int)
	go func() {
		i:=0
		for {
			time.Sleep(time.Duration(rand.Intn(500))*time.Millisecond)
			out <-i
			i++
		}
	}()
	return  out
}

func worker(id int,c chan int)  {
	for n:=range c{
		time.Sleep(800*time.Millisecond)
		fmt.Printf("Worker %d received %d\n",id,n)
	}
}

func createWork(id int) chan<-int{
	c:=make(chan int)
	go worker(id,c)
	return c
}

