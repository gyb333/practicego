package engine

import (
	"../config"
	"../parser"
	"../persist"
	"../scheduler"

	"../worker"
)



var url = config.Url

var cityUrl = config.CityUrl

var profileUrl = config.ProfileUrl
var hosts =[]int{
	9000,9001,9002,9003,9004,
}

func main() {
	//Simple()
	ConCurrent()
}

func ConCurrent() {
	itemChan, err := persist.ItemPrint(config.ItemSaverPort)
	if err !=nil{
		return
	}
	ConCurrentEngine{
		Scheduler:   &scheduler.QueuedScheduler{},
		WorkerCount: 10,
		ItemChan:   itemChan ,
		//ItemChan:  persist.SimpleItemPrint(),// itemChan ,
		//Worker: engine.Worker{},
		Worker:worker.RpcWorker{
			WorkChan:worker.CreatWorkerPool(hosts),
		},
	}.Run(
		Request{
			Url:    url,
			Parser: parser.CityListParser{},
		},
		//engine.Request{
		//	Url:    cityUrl,
		//	Parser: parser.CityParser{},
		//},
		//engine.Request{
		//	Url:    profileUrl,
		//	Parser: parser.ProfileParser{},
		//},
	)
}

func Simple() {
	SimpleEngine{
		ItemChan: persist.SimpleItemPrint(),
	}.Run(
		Request{
			Url:    url,
			Parser: parser.CityListParser{},
		},
		//engine.Request{
		//	Url:    cityUrl,
		//	Parser: parser.CityParser{},
		//},
		//engine.Request{
		//	Url:    profileUrl,
		//	Parser: parser.ProfileParser{},
		//},
	)
}
