package persist

import (
	"../engine"
	"log"
)

func SimpleItemPrint() chan engine.Item {
	out := make(chan engine.Item)
	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("Item Saver: Got %d  item : %v", itemCount, item)
			itemCount++
		}
	}()
	return out
}
