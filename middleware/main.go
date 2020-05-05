package main

import (
	"fmt"
	"time"
	"runtime"
)

func main()  {
	//selectDefault()
	defer func() {
		fmt.Println(recover())
	}()
	panic(1)
}




func selectDefault() {
	a := make(chan int, 1)
	b := make(chan int, 1)
	c := make(chan int, 1)
	//先执行从a队列收到的消息，其次执行b队列收到的消息，最后执行c队列收到的消息，那么执行顺序就是a>b>c
	go func() {
		ticker := time.NewTicker(1*time.Second)
		for{
			select {
				case <-ticker.C:
					a<-1
					b<-2
					c<-3
				}
			}
	}()
	for{
		select {
		case <-a:
			fmt.Println("recv from a")
		default:
			select {
			case <-b:
				fmt.Println("recv from b")
			default:
				select {
				case <-c:
					fmt.Println("recv from c")
				default:
					runtime.Gosched()
				}
			}
		}
	}

}
