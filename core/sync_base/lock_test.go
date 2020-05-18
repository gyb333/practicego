package sync_base

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCondLock(t *testing.T) {
	var wg sync.WaitGroup
	cond := sync.NewCond(new(sync.Mutex))
	for i := 0; i < 3; i++ {
		go func(i int) {
			fmt.Println("协程", i, "启动。。。")
			wg.Add(1)
			defer wg.Done()
			cond.L.Lock()
			fmt.Println("协程", i, "加锁。。。")
			cond.Wait()
			fmt.Println("协程", i, "解锁。。。")
			cond.L.Unlock()
		}(i)
	}
	time.Sleep(time.Second)
	fmt.Println("主协程发送信号量。。。")
	cond.Signal()
	time.Sleep(time.Second)
	fmt.Println("主协程发送信号量。。。")
	cond.Signal()
	time.Sleep(time.Second)
	fmt.Println("主协程发送信号量。。。")
	cond.Signal()
	wg.Wait()
}

func TestOnce(t *testing.T) {
	var once sync.Once
	var wg sync.WaitGroup
	onceFunc := func() {
		fmt.Println("法师爱你们哟~")
	}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			once.Do(onceFunc) // 多次调用只执行一次
		}()
	}
	wg.Wait()
}

func TestMap(t *testing.T) {
	var scene sync.Map
	// 将键值对保存到sync.Map
	scene.Store("法师", 97)
	scene.Store("老郑", 100)
	scene.Store("兵哥", 200)
	// 从sync.Map中根据键取值
	fmt.Println(scene.Load("法师"))
	// 根据键删除对应的键值对
	scene.Delete("法师")
	// 遍历所有sync.Map中的键值对
	scene.Range(func(k, v interface{}) bool {
		fmt.Println("iterate:", k, v)
		return true
	})

}