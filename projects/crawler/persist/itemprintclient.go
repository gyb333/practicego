package persist

import (
	"log"
	"../../utils"
	"../engine"
)

func ItemPrint(host string) (chan engine.Item, error)  {
	client, err := utils.NewClient(host)
	if err != nil {
		return nil, err
	}
	out := make(chan engine.Item)
	go func() {
		//itemCount := 0
		for{
			item := <-out
			//log.Printf("Item Saver: Got %d  item : %v", itemCount, item)
			//itemCount++
			//调用Rpc 来保存item
			result := ""
			err := client.Call("ItemPrintService.Save", item, &result)
			if err != nil {
				log.Printf("ItemPrintService :error saving item %v : %v ", item, err)
			}

		}
	}()
	return out,nil
}

