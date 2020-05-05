package engine

type Item struct {
	Url string
	Type string
	Id string
	Payload interface{}
}

type ParseResult struct {
	Requests []Request
	Items []Item
}

type Request struct {
	Url string
	Parser Parser
}

type Parser interface {
	Parse([]byte) ParseResult
	Serialize() (string, interface{})
}

type Enginer interface {
	Run(seeds ...Request)
}


//接口
type Scheduler interface {
	Submit(request Request)
	ReadyNotifier
	Run()
	WorkerChan() chan Request
}

//接口
type ReadyNotifier interface {
	WorkerReady(chan Request)
}

type RequestFetcher interface {
	FetchRequest(Request) (ParseResult,error)
}