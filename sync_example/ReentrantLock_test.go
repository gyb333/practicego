package sync_example_test

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
)

/*
Java中的synchronized关键词以及LinkedBlockingDequeu中用到的ReentrantLock，都是可重入的。
Golang中的sync.Mutex是不可重入的。表示锁不能递归使用，则会死锁
Golang的核心开发者认为可重入锁是不好的设计，所以不提供
 */

type ReentrantLock struct {
	*sync.Mutex
	*sync.Cond
	owner int
	holdCount int
}
func NewReentrantLock() sync.Locker {
	rl := &ReentrantLock{}
	rl.Mutex= new(sync.Mutex)
	rl.Cond = sync.NewCond(rl.Mutex)
	return rl
}

func GetGoroutineId() int {
	defer func()  {
		if err := recover(); err != nil {
			fmt.Printf("panic recover:panic info:%v", err)     }
	}()

	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

func (rl *ReentrantLock) Lock() {
	me := GetGoroutineId()
	rl.Lock()
	defer rl.Unlock()

	if rl.owner == me {
		rl.holdCount++
		return
	}
	for rl.holdCount != 0 {
		rl.Wait()
	}
	rl.owner = me
	rl.holdCount = 1
}

func (rl *ReentrantLock) Unlock() {
	rl.Lock()
	defer rl.Unlock()

	if rl.holdCount == 0 || rl.owner != GetGoroutineId() {
		panic("illegalMonitorStateError")
	}
	rl.holdCount--
	if rl.holdCount == 0 {
		rl.Signal()
	}
}




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

}