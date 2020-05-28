package test

import (
	. ".."
	"fmt"
	"sync"
	"testing"
)


func TestCyclicBarrier(t *testing.T) {
	N:=4
	barrier := NewCyclicBarrier(N)
	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("A")
			barrier.BarrierWait()
			fmt.Println("B")
			barrier.BarrierWait()
			fmt.Println("C")
		}()
	}
	wg.Wait()
}


/*
// 抽象一个栅栏
type Barrier interface {
	Wait ()
}
// 创建栅栏对象
func NewBarrier (n int) Barrier {
}
// 栅栏的实现类
type barrier struct {
}
// 测试代码
func main () {
	// 创建栅栏对象
	b := NewBarrier(10)
	// 达到的效果：前9个协程调用Wait()阻塞，第10个调用后10个协程全部唤醒
	for i := 0; i < 10; i++ {
		go b.Wait()
	}
}
 */

func TestCyclicBarrierWait(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(10)
	// 创建栅栏对象
	b := NewCyclicBarrier(10)
	// 达到的效果：前9个协程调用Wait()阻塞，第10个调用后10个协程全部唤醒
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			fmt.Printf("调用Wait阻塞,协程序号%d......\n",i)
			b.BarrierWait()
			fmt.Printf("执行Wait完成,协程序号%d......\n",i)
		}(i)
	}
	wg.Wait()
}
