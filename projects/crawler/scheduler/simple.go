package scheduler

import "../engine"

/**
 * 非队列，公用一个channel
 * 全都在抢一个channel，不可控制
 */

type SimpleScheduler struct {
	workerChan chan engine.Request
}

func (s *SimpleScheduler) Submit(request engine.Request) {
	go func() {
		s.workerChan<-request
	}()
}

func (s *SimpleScheduler) WorkerReady(chan engine.Request) {

}

func (s *SimpleScheduler) Run() {
	s.workerChan=make(chan engine.Request)
}

func (s *SimpleScheduler) WorkerChan() chan engine.Request {
	return  s.workerChan
}




