package sync_example_test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

type RWLocker struct {
	sync.RWMutex
	data interface{}//共享数据，只能有一个线程能写该数据，但可以有多个线程同时读该数据。
}

func (l *RWLocker) Get()  {
	l.RLock()
	fmt.Println("获取data数据,之前执行");
	//设置随机种子
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(1000)))
	fmt.Printf("获取data数据：%v\n",l.data)
	l.RUnlock()
}

func (l *RWLocker) Set(data interface{})  {
	l.Lock()
	fmt.Printf("设置数据前执行,%d\n", l.data)
	//设置随机种子
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(1000)))
	l.data=data
	fmt.Printf("设置数据后执行,%d\n",l.data)
	l.Unlock()
}

func TestRWLocker(t *testing.T) {
	wg :=sync.WaitGroup{ }
	l:=RWLocker{}
	N:=20
	wg.Add(N)
	for i:=0;i<N;i++{
		go func() {
			defer wg.Done()
			for {
				l.Get()
			}
		}()
		go func() {
			defer wg.Done()
			for{
				rand.Seed(time.Now().UnixNano())
				l.Set(rand.Intn(1000))
			}

		}()
	}
	wg.Wait()
}
