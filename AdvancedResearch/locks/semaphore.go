package locks

import "sync"

type semaphore struct {
	refs int

	max int
	value int
	sync.Cond
}

func NewSemaphore(max int) *semaphore  {
	return &semaphore{
		max: max,
		Cond: sync.Cond{L:new(sync.Mutex)},
	}
}

func (s *semaphore) Values() int {
	s.L.Lock()
	defer s.L.Unlock()
	return s.value
}

func (s *semaphore) Acquire()  {
	s.L.Lock()
	defer s.L.Unlock()
	for{
		if s.value+1<=s.max{
			s.value++
			return
		}
		s.Wait()
	}
	panic("unexpected branch")
}

func (s *semaphore) Release(){
	s.L.Lock()
	defer s.L.Unlock()
	s.value--
	if s.value<0{
		panic("semaphore Release without Acquire")
	}
	s.Signal()
}