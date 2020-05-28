package _map
import (
	"sync"
)

type RWMap struct {
	sync.RWMutex
	KeyValues map[string] interface{}
}

func NewRWMap()  *RWMap{
	return &RWMap{
		KeyValues:make(map[string]interface{}),
	}
}

func (rw *RWMap) Get(key string)(value interface{},ok bool){
	rw.RLock()
	result,ok :=rw.KeyValues[key]
	rw.RUnlock()
	return result,ok
}

func (rw *RWMap) Delete(key string){
	rw.Lock()
	delete(rw.KeyValues,key)
	rw.Unlock()
}

func (rw *RWMap) Set(key string ,value interface{}){
	rw.Lock()
	rw.KeyValues[key]=value
	rw.Unlock()
}