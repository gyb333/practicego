package basic

import (
	"runtime"
	"fmt"
	"os"
)

/*
	goroutine 可能的切换点
	I/O Select
	channel
	等待锁
	函数调用(有时)
	runtime.Gosched() time.Sleep()
 */
func GoRoutineMain()  {
	var a [10]int
	for i:=0;i<10;i++{
		go func(i int){
			for {
			a[i]++
			runtime.Gosched()			//需要主动交出控制权
			//fmt.Printf("Hello from goroutine %d\n",i)	//io操作会交出控制权
			}
		}(i)
	}
	//time.Sleep(time.Microsecond)
	var str string
	fmt.Fscanln(os.Stdin,&str)
	fmt.Println(a)
}
