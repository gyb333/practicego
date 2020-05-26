package sync_example_test

import (
	"fmt"
	"sync"
	"testing"
)

/*
N个线程，按顺序交替执行打印,要求结果如下：
*/

type Conditions struct {
	//mtx sync.Mutex		//使用一个可重入的锁ReentrantLock
	cs []sync.Cond
	index int
	length int
}

func NewConditions(length int)  *Conditions{
	cs:=make([]sync.Cond,length)
	for i:=range cs{
		cs[i]=sync.Cond{L: new(sync.Mutex)}
	}
	return &Conditions{cs:cs,length: length}
}

func (c *Conditions) DoFunc(index int,Func func()){
	c.cs[index].L.Lock()
	defer c.cs[index].L.Unlock()
	for c.index!=index{
		c.cs[index].Wait()
	}
	Func()
	c.index=(index+1)%c.length
	c.cs[(index+1)%c.length].Signal()
}

const (
	  COUNTS  =1000
	  GOS     =10
)
func TestConditions(t *testing.T) {
	wg:=sync.WaitGroup{}
	wg.Add(GOS)
	c:=NewConditions(GOS)
	count:=0
	for i:=0;i<c.length;i++{
		go func(i int) {
			defer wg.Done()
			for count<=COUNTS{
				c.DoFunc(i, func() {
					count++
					if count<=COUNTS{
						fmt.Printf("协程%2d,执行数据%4d\n",i+1,count)
					}
				})
			}
		}(i)
	}
	wg.Wait()
}