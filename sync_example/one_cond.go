package sync_example

import "sync"

type OneCondition struct {
	sync.Cond
	bShouldSub bool
}

func NewOneCondition()  *OneCondition{
	return &OneCondition{Cond:sync.Cond{
		L:new(sync.Mutex),
	},
	}
}

func (b *OneCondition) DoNextFunc(Func func())  {
	b.L.Lock()
	defer b.L.Unlock()
	for !b.bShouldSub{
		b.Wait()
	}
	Func()
	b.bShouldSub = false
	b.Signal()
}
func (b *OneCondition) DoFrontFunc(Func func())  {
	b.L.Lock()
	defer b.L.Unlock()
	for b.bShouldSub{
		b.Wait()
	}
	Func()
	b.bShouldSub = true
	b.Signal()
}

func (b *OneCondition) DoFunc(isNext bool,Func func())  {
	b.L.Lock()
	defer b.L.Unlock()
	if isNext{
		for !b.bShouldSub{
			b.Wait()
		}
	}else{
		for b.bShouldSub{
			b.Wait()
		}
	}

	Func()
	if isNext{
		b.bShouldSub = false
	}else {
		b.bShouldSub = true
	}
	b.Signal()
}

