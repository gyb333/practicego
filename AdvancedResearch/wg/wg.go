package wg

import (
	"sync"
	"time"
)

type WaitGroupWrapper struct {
	wg sync.WaitGroup
}

func (this *WaitGroupWrapper) Wrap(f func()) {
	this.wg.Add(1)
	go func() {
		defer this.wg.Done()
		f()
	}()
}

func (this *WaitGroupWrapper) Wait() {
	this.wg.Wait()
}

// WaitTimeout is same as Wait except that it accepts timeout arguement.
// FIXME if timeout triggered, there will be goroutine leak.
func (this *WaitGroupWrapper) WaitTimeout(timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		this.wg.Wait()
	}()

	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
