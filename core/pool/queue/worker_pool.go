package queue

import (
	. "../../pool"
)

type WorkerPool struct {
	Scheduler Scheduler
	WorkerCount int
}

func NewWorkerPool(workerCount int) *WorkerPool{
	return &WorkerPool{
		WorkerCount: workerCount,
		Scheduler:&QueueScheduler{},
	}
}

func (e WorkerPool) Run( ){
	e.Scheduler.Run()
	for i := 0; i < e.WorkerCount; i++ {
		e.CreateWorker(e.Scheduler.WorkerChan(), e.Scheduler)
	}
}

func (e WorkerPool)CreateWorker(in chan Job,ready  ReadyNotifier) {
	go func() {
		for {
			//需要让scheduler知道已经就绪了
			ready.WorkerReady(in)
			job := <-in
			job.Do()
		}
	}()
}

