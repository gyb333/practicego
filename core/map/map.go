package _map

import (
	"reflect"
	"sync"
)

/*
key中不能包含不可比较的值，比如 slice, map, and function。
而我们的key是用户自定义的对象，没办法进行约束。于是借鉴Java的IdentityHashMap的思路，
将key转换成对象的指针地址，实际上map中保存的是key对象的指针地址。
 */


type SyncIdentityMap struct {
	sync.RWMutex
	m map[uintptr]interface{}
}

func (this *SyncIdentityMap) Get(key interface{}) interface{} {
	this.RLock()
	keyPtr := genKey(key)
	value := this.m[keyPtr]
	this.RUnlock()
	return value
}

func genKey(key interface{}) uintptr {
	keyValue := reflect.ValueOf(key)
	return keyValue.Pointer()
}