package test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestExchangerData(t *testing.T)  {
	wg:=sync.WaitGroup{}
	wg.Add(2)
	cdata:=make(chan int)
	go func() {
		defer wg.Done()
		for data:=range cdata{
			fmt.Printf("读协程获取数据%d\n",data)
		}
	}()
	go func() {
		defer wg.Done()
		for{
			data:=rand.Intn(100)
			select {
				case cdata<-data:
					fmt.Printf("写协程设置数据%d\n",data)
				case <-time.After(time.Nanosecond):
					close(cdata)
					return
			}
		}
	}()

	wg.Wait()
}
