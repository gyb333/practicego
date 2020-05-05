package main

import (
	"utils"
	"log"
		"crawler/engine"
	"crawler/config"
)

type ItemPrintService struct {

}
var count=0
func (s ItemPrintService) Save(item engine.Item, result *string) error {
	count++
	log.Println(count,item)
	*result = "ok"
	return nil
}
func main() {
	log.Fatal(utils.ServeRpc(config.ItemSaverPort,ItemPrintService{}))
}
