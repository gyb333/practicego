package pool

type Job interface {
	Do()
}



type Scheduler interface {
	Submit(job Job)
	ReadyNotifier
	Run()
	WorkerChan() chan Job
}

//接口
type ReadyNotifier interface {
	WorkerReady(chan Job)
}

