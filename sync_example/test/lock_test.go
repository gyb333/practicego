package test_test

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

type Locker struct {
	sync.Mutex
}

func NewLocker()  *Locker{
	return &Locker{sync.Mutex{},}
}

func (l *Locker) print(name string)  {
	l.Lock()
	defer func() {
		l.Unlock()
		runtime.Gosched()		//主动让出资源
	}()
	for _,v:=range name{
		fmt.Printf("%c",v)
	}
	fmt.Println()

}


func TestLocker(t *testing.T) {
	wg :=sync.WaitGroup{}
	wg.Add(2)
	l :=NewLocker()
	go func() {
		defer wg.Done()
		t:=time.NewTicker(10*time.Millisecond)
		for {
			select {
			case <-t.C:
				return
			default:
				l.print("ggggggggggggggg")
			}
		}
	}()
	go func() {
		defer wg.Done()
		t:=time.NewTicker(10*time.Millisecond)
		for {
			select {
			case <-t.C:
				return
			default:
				l.print("ooooooooooooooo")
			}
		}
	}()
	wg.Wait()
}
