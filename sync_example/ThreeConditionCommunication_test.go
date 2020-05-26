package sync_example_test

import (
	"fmt"
	"sync"
	"testing"
)

type ThreeCondition struct {
	mtx sync.Mutex		//使用一个可重入的锁ReentrantLock
	cs [3]sync.Cond
	Sub int
}

func NewThreeCondition()  *ThreeCondition{
	m :=sync.Mutex{}
	return &ThreeCondition{
		mtx:m,
		cs :[3]sync.Cond{
		sync.Cond{L: new(sync.Mutex)},	//golang 目前没有可重入锁，用多个mutex代替
		sync.Cond{L: new(sync.Mutex)},
		sync.Cond{L: new(sync.Mutex)},
		},
	}
}

func (t *ThreeCondition) sub2(i int)  {
	t.cs[1].L.Lock()
	defer t.cs[1].L.Unlock()
	for t.Sub!=1{
		t.cs[1].Wait()
	}
	for j:=1;j<=10;j++ {
		fmt.Printf("sub2 thread sequence of %d,loop of %d\n" ,j, i);
	}
	t.Sub=2
	t.cs[2].Signal()
}
func (t *ThreeCondition) sub3(i int)  {
	t.cs[2].L.Lock()
	defer t.cs[2].L.Unlock()
	for t.Sub!=2{
		t.cs[2].Wait()
	}
	for j:=1;j<=5;j++ {
		fmt.Printf("sub3 thread sequence of %d,loop of %d\n" ,j, i);
	}
	t.Sub=0
	t.cs[0].Signal()
}

func (t *ThreeCondition) main(i int)  {
	t.cs[0].L.Lock()
	defer t.cs[0].L.Unlock()
	for t.Sub!=0{
		t.cs[0].Wait()
	}
	for j:=1;j<=6;j++ {
		fmt.Printf("main thread sequence of %d,loop of %d\n" ,j, i);
	}
	t.Sub = 1
	t.cs[1].Signal()
}

func TestThreeCondition(t *testing.T) {
	business:=NewThreeCondition()
	go func() {
		for i:=1;i<=50;i++ {
			business.sub2(i);
		}
	}()

	go func() {
		for i:=1;i<=50;i++ {
			business.sub3(i);
		}
	}()

	for i:=1;i<=50;i++ {
		business.main(i)
	}
}