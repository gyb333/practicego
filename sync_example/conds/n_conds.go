package conds

import "sync"

/*
Golang的 sync.Cond 只有Wait，没有如Java中的Condition的超时等待方法await(long time, TimeUnit unit)。
没法实现LinkBlockingDeque的 pollFirst(long timeout, TimeUnit unit) 这样的方法。
 */
type NConditions struct {
	//mtx sync.Mutex		//使用一个可重入的锁ReentrantLock
	cs []sync.Cond
	index int
	Length int
}

func NewNConditions(length int)  *NConditions {
	cs:=make([]sync.Cond,length)
	for i:=range cs{
		cs[i]=sync.Cond{L: new(sync.Mutex)}
	}
	return &NConditions{cs: cs,Length: length}
}

func (c *NConditions) DoFunc(index int,Func func()){
	c.cs[index].L.Lock()
	defer c.cs[index].L.Unlock()
	for c.index!=index{
		c.cs[index].Wait()
	}
	Func()
	c.index=(index+1)%c.Length
	c.cs[(index+1)%c.Length].Signal()
}
