package sync_example

import "sync"

type CyclicBarrier struct {
	curCnt int
	maxCnt int
	*sync.Cond
}

func NewCyclicBarrier(maxCnt int) *CyclicBarrier {
	return &CyclicBarrier{
		curCnt: maxCnt,
		maxCnt: maxCnt,
		Cond:&sync.Cond{L: new(sync.Mutex)},
	}
}

func (cb *CyclicBarrier)BarrierWait()  {
	cb.L.Lock()
	defer cb.L.Unlock()
	if cb.curCnt--;cb.curCnt>0{
		cb.Wait()
	}else {
		cb.Broadcast()
		cb.curCnt=cb.maxCnt
	}
}
