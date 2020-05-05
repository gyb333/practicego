package worker

import (
	"../../utils"
	"../config"
	"../engine"
	"fmt"
	"log"
	"net/rpc"
)

type RpcWorker struct {
	WorkChan chan *rpc.Client
}

func (w RpcWorker) FetchRequest(request engine.Request) (engine.ParseResult, error) {
	//client, err := utils.NewClient(fmt.Sprintf(":%d", config.WorkerPort0))
	//if err != nil {
	//	return engine.ParseResult{}, err
	//}
	client :=<-w.WorkChan

	sReq := SerializeRequest(request)
	var sResult ParseResult
	err := client.Call(config.CrawlServiceRpc, sReq, &sResult)
	if err != nil {
		return engine.ParseResult{}, err
	}
	return DeserializeResult(sResult), nil

}

func CreatWorkerPool(hosts []int) chan *rpc.Client{
	var clients []*rpc.Client
	for _,host:=range hosts{
		client,err:=utils.NewClient(fmt.Sprintf(":%d", host))
		if err!=nil{
			log.Printf("Error connecting to %s %v\n",host,err)
		}else{
			log.Printf("Success Connecting to %d\n",host)
			clients=append(clients,client)
		}
	}
	out :=make(chan *rpc.Client)
	go func(){
		for {
			for _, client := range clients {
				out <- client
			}
		}
	}()
	return out
}
