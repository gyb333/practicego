package worker

import (
	"../engine"
	"../fetcher"
	"fmt"
	"log"
)

type FetchWorker struct {

}

func (FetchWorker)FetchRequest(request engine.Request) (engine.ParseResult,error) {
	return  FetchWork(request)

}


func FetchWork(request engine.Request) (engine.ParseResult,error) {
	log.Printf("Fetching %s",request.Url)
	body,err:=fetcher.Fetch(request.Url)
	if err!=nil{
		return engine.ParseResult{},err
	}
	if len(body)<=0{
		err= fmt.Errorf("Body contents %d URL:%s\n",len(body),request.Url)
		return engine.ParseResult{},err
	}
	//log.Printf("%s",string(body))
	return request.Parser.Parse(body),nil

}