package sizedwaitgroup

import (
	"context"
	"math"
	"sync"
)

type SizeWaitGroup struct {
	Size int
	current chan struct{}
	wg sync.WaitGroup
}


func New(limit int) SizeWaitGroup{
	size :=math.MaxInt32
	if limit>0{
		size =limit
	}
	return SizeWaitGroup{Size: size,
		current: make(chan struct{},size),
		wg:sync.WaitGroup{},
	}
}

func (swg *SizeWaitGroup) Add(){
	swg.AddWithContext(context.Background())
}

func (swg *SizeWaitGroup) AddWithContext(ctx context.Context) error{
	select {
	case <-ctx.Done():
		return ctx.Err()
	case swg.current<- struct{}{}:
		break
	}
	swg.wg.Add(1)
	return nil
}

func (swg *SizeWaitGroup) Done()  {
	<-swg.current
	swg.wg.Done()
}

func (swg *SizeWaitGroup) Wait()  {
	swg.wg.Wait()
}