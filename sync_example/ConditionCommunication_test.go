package sync_example_test

import (
	"fmt"
	"sync"
	"testing"
)

type Business struct {
	sync.Cond
	bShouldSub bool
}

func NewBusiness()  *Business{
	return &Business{Cond:sync.Cond{
			L:new(sync.Mutex),
		},
	}
}

func (b *Business) sub(i int)  {
	b.L.Lock()
	defer b.L.Unlock()
	for !b.bShouldSub{
		b.Wait()
	}
	for j:=1;j<=10;j++ {
		fmt.Printf("sub thread sequence of %d,loop of %d\n" ,j, i);
	}
	b.bShouldSub = false
	b.Signal()
}
func (b *Business) main(i int)  {
	b.L.Lock()
	defer b.L.Unlock()
	for b.bShouldSub{
		b.Wait()
	}
	for j:=1;j<=20;j++ {
		fmt.Printf("main thread sequence of %d,loop of %d\n" ,j, i);
	}
	b.bShouldSub = true
	b.Signal()
}



func TestCondition(t *testing.T) {
	business := NewBusiness()
	go func() {
		for i:=1;i<=50;i++ {
			business.sub(i);
		}
	}()
	for i:=1;i<=50;i++ {
		business.main(i);
	}
}
