package sync_example

import "sync"


type OneBarrier struct {
	c      int
	n      int
	m      sync.Mutex
	before chan struct{}
	after  chan struct{}
}

func NewOneBarrier(n int) *OneBarrier {
	b := OneBarrier{
		n:      n,
		before: make(chan struct{}, 1),
		after:  make(chan struct{}, 1),
	}
	// close 1st gate
	b.after <- struct{}{}
	return &b
}

func (b *OneBarrier) Before() {
	b.m.Lock()
	b.c += 1
	if b.c == b.n {
		// close 2nd gate
		<-b.after
		// open 1st gate
		b.before <- struct{}{}
	}
	b.m.Unlock()
	<-b.before
	b.before <- struct{}{}
}

func (b *OneBarrier) After() {
	b.m.Lock()
	b.c -= 1
	if b.c == 0 {
		// close 1st gate
		<-b.before
		// open 2st gate
		b.after <- struct{}{}
	}
	b.m.Unlock()
	<-b.after
	b.after <- struct{}{}
}

