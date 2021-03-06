package simple

import (
	. "../../pool"
	"runtime"
)

type Worker struct {
	JobQueue chan Job
}

func NewWorker() Worker  {
	return Worker{JobQueue: make(chan Job)}
}

func (w Worker) Run(work chan chan Job)  {
	go func() {
		for{
			select {
				case work<-  w.JobQueue:
				case job:=<-w.JobQueue:
					job.Do()
			}
		}
	}()
}


type WorkerPool struct {
	MaxWorker int
	JobQueue chan Job
	WorkerQueue chan chan Job
}

func NewWorkerPool(maxWorker int) *WorkerPool{
	return &WorkerPool{
		MaxWorker: maxWorker,
		JobQueue: make(chan Job),
		WorkerQueue: make(chan chan Job),
	}
}

func (wp *WorkerPool) Run(){
	for i:=0;i<wp.MaxWorker;i++{
		worker:=NewWorker()
		worker.Run(wp.WorkerQueue)
	}
	go func() {
		for{
			select {
				case job :=<-wp.JobQueue:
					worker:=<-wp.WorkerQueue
					worker<-job
				default:
					runtime.Gosched()
			}
		}
	}()
}