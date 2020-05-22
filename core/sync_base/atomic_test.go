package sync_base_test

import (
	"fmt"
	"sync/atomic"
	"testing"
)

func TestAtomic(t *testing.T) {
	var counter int64 =23
	atomic.AddInt64(&counter,-3)
	fmt.Println(counter)

	for{
		v :=atomic.LoadInt64(&counter)	//原子读取
		if atomic.CompareAndSwapInt64(&counter,v,v+10){
			break
		}
	}
	fmt.Println(counter)

	var atomicValue atomic.Value	//atomic.Value类型的变量一旦被声明，就不应该被复制到其他地方
	atomicValue.Store([]int{1,2,3,4,5})
	anotherStore(atomicValue)			//拷贝了一个副本值
	fmt.Println("main: ",atomicValue)

}
func anotherStore(Atomicvalue atomic.Value)  {
	Atomicvalue.Store([]int{6,7,8,9,10})
	fmt.Println("anotherStore: ",Atomicvalue)
}