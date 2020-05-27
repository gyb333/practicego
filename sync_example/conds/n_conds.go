package conds

import "sync"

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
