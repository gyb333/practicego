package engine

import (
	"log"
	)

type SimpleEngine struct {
	ItemChan chan Item
	Worker RequestFetcher
}

func (s SimpleEngine) Run(seeds ...Request){
	var requests []Request
	requests = append(requests,seeds...)
	for len(requests)>0{
		request:=requests[0]
		requests=requests[1:]
		//log.Printf("%s",string(body))
		parseResult,err:=s.Worker.FetchRequest(request)
		if err!=nil{
			log.Printf("Fetcher: error fetching url %s %v",request.Url,err)
			continue
		}
		//requests=append(requests,parseResult.Requests...)
		for _,item:=range parseResult.Items{
			//log.Printf("Got item %d %v",i,item)
			s.ItemChan<-item
		}
	}
}
