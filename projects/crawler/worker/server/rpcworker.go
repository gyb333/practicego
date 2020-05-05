package main

import (
	"../../../utils"
	"../../worker"
	"flag"
	. "fmt"
	"log"
)

type CrawlService struct {
}



func (CrawlService) Process(req worker.Request, result *worker.ParseResult) error {
	engineReq, err := worker.DeserializeRequest(req)
	if err != nil {
		return err
	}
	engineRes, err := worker.FetchWork(engineReq)
	if err != nil {
		return err
	}

	*result = worker.SerializeParseResult(engineRes)
	return nil
}


var port = flag.Int("port", 9000, "the port for me to listen on")
//go run rpcworker.go --port=9000
func main() {
	flag.Parse()
	if *port == 0 {
		log.Println("must specify a port ... ")
		return
	}

	log.Fatal(utils.ServeRpc(Sprintf(":%d", *port), CrawlService{}))


}
