package persist

import (
	"log"

	"../engine"

	"../../utils"
)

func ItemSaver(host string) (chan engine.Item, error) {
	client, err := utils.NewClient(host)
	if err != nil {
		return nil, err
	}

	out := make(chan engine.Item)
	go func() {
		for {
			item := <-out
			//调用Rpc 来保存item
			result := ""
			err = client.Call("ItemSaveService.Save", item, &result)
			if err != nil {
				log.Printf("ItemSaveService :error saving item %v : %v ", item, err)
			}
		}
	}()
	return out, nil

}
