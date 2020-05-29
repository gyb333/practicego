package test_test

import (
	. ".."
	"fmt"
	"sync"
	"testing"
)




type LockStruct struct {
	sync.Locker
	name string
	id int
}

func (s *LockStruct) setName(name string) {
	s.Lock()
	defer s.Unlock()
	s.name = name
}

func (s *LockStruct) setId(id int) {
	s.Lock()
	defer s.Unlock()
	s.id = id
}

func (s LockStruct) PrintName() {
	s.Lock()
	defer s.Unlock()
	s.setName("goroutine id : ")
	s.setId(GetGoroutineId())
	fmt.Println(s.name, s.id)
}


func TestReentrantLock(t *testing.T) {
	fmt.Println("reentrant lock single goroutine test start")
	ls := &LockStruct{Locker: NewReentrantLock()}
	ls.PrintName()
	fmt.Println("reentrant lock single goroutine test end")

	wg:=sync.WaitGroup{}
	N:=10
	wg.Add(N)
	fmt.Println("reentrant lock multi goroutine test start")
	for i:=0;i<N;i++{
		go func() {
			defer wg.Done()
			ls := &LockStruct{Locker: NewReentrantLock()}
			ls.PrintName()
		}()
	}
	wg.Wait()
	fmt.Println("reentrant lock multi goroutine test end")
}