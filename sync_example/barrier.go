package sync_example

import "sync"

/*
使用了容量和 worker 数量 n 相等的 buffered channel。
不再让 worker一个接一个地依次通过，而是在channel中放入n个元素来使所有的 worker一次性通过
 */

type Barrier struct {
	c      int
	n      int
	m      sync.Mutex
	before chan struct{}
	after  chan struct{}
}

func NewBarrier(n int) *Barrier {
	b := Barrier{
		n:      n,
		before: make(chan struct{}, n),
		after:  make(chan struct{}, n),
	}
	return &b
}

func (b *Barrier) Before() {
	b.m.Lock()
	b.c += 1
	if b.c == b.n {
		// open 2nd gate
		for i := 0; i < b.n; i++ {
			b.before <- struct{}{}
		}
	}
	b.m.Unlock()
	<-b.before
}

func (b *Barrier) After() {
	b.m.Lock()
	b.c -= 1
	if b.c == 0 {
		// open 1st gate
		for i := 0; i < b.n; i++ {
			b.after <- struct{}{}
		}
	}
	b.m.Unlock()
	<-b.after
}