package locks

import (
	"sync"
)

type blockingKeyCountLimit struct {
	sync.RWMutex
	current map[string]*semaphore
	limit int
}

func NewBlockingKeyCountLimit(n int) *blockingKeyCountLimit{
	return &blockingKeyCountLimit{
		current: make(map[string]*semaphore),
		limit: n,
	}
}

func (bl *blockingKeyCountLimit) Values() int{
	bl.RLock()
	defer bl.RUnlock()
	all:=0
	for _,v:=range bl.current{
		all+=v.Values()
	}
	return all
}

func (bl *blockingKeyCountLimit) Acquire(key []byte) {
	strKey:=string(key)
	bl.Lock()
	k,ok:=bl.current[strKey];
	if !ok{
		k=NewSemaphore(bl.limit)
		bl.current[strKey]=k
	}
	k.refs++
	bl.Unlock()
	k.Acquire()
}

func (bl *blockingKeyCountLimit) Release(key []byte)  {
	strKey :=string(key)
	bl.Lock()
	k,ok:=bl.current[strKey]
	if !ok{
		panic("key not in map. Possible reason: Release without Acquire.")
	}
	k.refs--
	if k.refs<0{
		panic("internal error: refs < 0")
	}
	if k.refs==0{
		delete(bl.current,strKey)
	}
	bl.Unlock()
	k.Release()
}