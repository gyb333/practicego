package test

import (
	. "SyncExample"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)


func init() {
	rand.Seed(time.Now().Unix())
}

type counter struct {
	c int
	sync.Mutex
}

func (c *counter) Incr() {
	c.Lock()
	c.c += 1
	c.Unlock()
}

func (c *counter) Get() (res int) {
	c.Lock()
	res = c.c
	c.Unlock()
	return
}

/*
只有 n, 2*n 和 3*n 的数字被打印（因为每个 worker 循环 3 次）
每次打印的数字不会比之前的数字小
每个数字会被打印 n 次

如果有 3 个 worker，则期望输出如下：
3
3
3
6
6
6
9
9
9
2 个 worker 的期望输出：
2
2
4
4
6
6
2 个 worker 不合法的输出可能会是这样：
2
4
2
4
6
6
*/

var workers = 100
var Count=100
func TestOneBarrier(t *testing.T) {
	var wg sync.WaitGroup
	br := NewOneBarrier(workers)
	c := counter{}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			for i := 0; i < Count; i++ {
				br.Before()
				c.Incr()
				br.After()
				//time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
				fmt.Printf("goroutine %d,获取数据：%d\n",j,c.Get())
			}
		}(i+1)
	}
	wg.Wait()
}

func TestBarrier(t *testing.T) {
	var wg sync.WaitGroup
	br := NewBarrier(workers)
	c := counter{}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			for i := 0; i < Count; i++ {
				br.Before()
				c.Incr()
				br.After()
				//time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
				fmt.Printf("goroutine %d,获取数据：%d\n",j,c.Get())
			}
		}(i+1)
	}
	wg.Wait()
}

func TestCBarrier(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(workers)
	br := NewCyclicBarrier(workers)
	c := counter{}
	for i := 0; i < workers; i++ {
		go func(j int) {
			defer wg.Done()
			for i := 0; i < Count; i++ {
				br.BarrierWait()
				c.Incr()
				br.BarrierWait()
				//time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
				fmt.Printf("goroutine %d,获取数据：%d\n",j,c.Get())
			}
		}(i+1)
	}
	wg.Wait()
}