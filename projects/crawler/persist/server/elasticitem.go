package main

import (
	"gopkg.in/olivere/elastic.v5"

	"utils"
	"log"
	"crawler/engine"

	"crawler/persist/elastics"
)



type ItemSaveService struct {
	Client *elastic.Client
	Index  string
}

func (s *ItemSaveService) Save(item engine.Item, result *string) error {
	err := elastics.Save(s.Client, s.Index, item)
	if err == nil {
		*result = "ok"
	}
	return err
}

func main() {
	//如果发生错误，Fatal()会强制退出。。
	log.Fatal(serveRpc("7788", "datint_profile"))
}

func serveRpc(host, index string) error {
	client, err := elastic.NewClient(
		elastic.SetSniff(false))
	if err != nil {
		return err
	}
	utils.ServeRpc(host, &ItemSaveService{
		Client: client,
		Index:  index,
	})
	return nil
}