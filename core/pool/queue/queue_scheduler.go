package queue


import (
	. "../../pool"
)

type QueueScheduler struct {
	jobChan chan Job
	workerChan  chan chan Job
}


func (s *QueueScheduler) WorkerChan() chan Job {
	return make(chan Job)
}

func (s *QueueScheduler) WorkerReady(w chan Job) {
	s.workerChan <- w
}

func (s *QueueScheduler) Submit(job Job) {
	s.jobChan <- job
}
func (s *QueueScheduler) Run() {
	s.workerChan = make(chan chan Job)
	s.jobChan = make(chan Job)
	go func() {
		var requestQ []Job
		var workerQ []chan Job
		for {
			var activeRequest Job
			var activeWorker chan Job
			if len(requestQ) > 0 && len(workerQ) > 0 {
				// worker和request同时有闲置的时，才给active...赋值，否则均为初始值nil
				activeWorker = workerQ[0]
				activeRequest = requestQ[0]
			}
			// 独立的两件事用select分别处理，它们发生的前后不定
			select {
			case r := <-s.jobChan:
				// 有request来到request队列，则加入
				requestQ = append(requestQ, r)
			case w := <-s.workerChan:
				// 有worker来到worker队列，则加入
				workerQ = append(workerQ, w)
			case activeWorker <- activeRequest:
				// 当activeWorker是初始值nil，则不会执行该select case
				// 当有worker，执行该case，则把已经使用的worker和request从队列里拿掉
				workerQ = workerQ[1:]
				requestQ = requestQ[1:]
			}
		}
	}()
}