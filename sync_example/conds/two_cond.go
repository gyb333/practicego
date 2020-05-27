package conds

import "sync"

type TwoCondition struct {
	frontCond sync.Cond
	nextCond  sync.Cond
	bFlag     bool
}

func NewTwoCondition()  *TwoCondition{
	return &TwoCondition{
		frontCond:sync.Cond{L:new(sync.Mutex)},
		nextCond: sync.Cond{L: new(sync.Mutex)},
	}
}

func (b *TwoCondition) DoNextFunc(Func func())  {
	b.nextCond.L.Lock()
	defer b.nextCond.L.Unlock()
	for !b.bFlag {
		b.nextCond.Wait()
	}
	Func()
	b.bFlag = false
	b.frontCond.Signal()
}
func (b *TwoCondition) DoFrontFunc(Func func())  {
	b.frontCond.L.Lock()
	defer b.frontCond.L.Unlock()
	for b.bFlag {
		b.frontCond.Wait()
	}
	Func()
	b.bFlag = true
	b.nextCond.Signal()
}
