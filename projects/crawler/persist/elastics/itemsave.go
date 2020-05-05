package elastics

import (
	"gopkg.in/olivere/elastic.v5"
	"crawler/engine"

	"errors"
	"context"
)





func Save(client *elastic.Client,index string, item engine.Item) (err error) {
	if item.Type == "" {
		return errors.New("must supply Type")
	}
	indexService := client.Index().Index(index).Type(item.Type).BodyJson(item)
	if item.Id != "" {
		indexService.Id(item.Id)
	}
	_, err = indexService.Do(context.Background())
	if err != nil {
		return err
	}
	//fmt.Printf("%+v",response)  //+v带有key
	return nil
}