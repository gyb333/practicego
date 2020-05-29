package test_test

import (
	. ".."
	"fmt"
	"sync"
	"testing"
)

/*
请编写2个线程，线程1顺序输出1，3，5，……, 99 等奇数，每个数 一 。
线程2顺序输出2，4，6……100等偶数，每个数 一 。
最终的结果要求是输出为 自然顺序：1, 2, 3, 4, ……99, 100。
*/



func TestChan(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	ch := make(chan int)

	go func() {
		wg.Done()
		for i := 1; i <= 50; i++ {
			println("协程g1:", <-ch) // 执行步骤1， 执行步骤5
			ch <- i + i            // 发生偶数
		}
	}()

	go func() {
		defer func() {
			wg.Done()
			close(ch)
		}()
		for i := 1; i <= 50; i++ {
			ch <- i + i - 1        //发送
			println("协程g2:", <-ch) //执行步骤4
		}
	}()
	wg.Wait()
}

func TestChan2(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	c1 := make(chan int)
	c2 := make(chan int)
	go func() {
		defer wg.Done()
		for j := 1; j <= 50; j++ {
			fmt.Printf("协程1,输出：%d\n", <-c1)
			c2 <- j + j
		}
	}()
	go func() {
		defer wg.Done()
		for j := 1; j <= 50; j++ {
			c1 <- j + j - 1
			fmt.Printf("协程2,输出：%d\n", <-c2)
		}
	}()
	wg.Wait()
	close(c1)
	close(c2)
}

func TestCond(t *testing.T) {
	wg:= sync.WaitGroup{}
	wg.Add(2)
	oneCond:=NewOneCondition()
	count:=0

	go func() {
		defer wg.Done()
		for count<100{
			oneCond.DoFunc(true,func() {
					count++
					fmt.Printf("go 协程 2,loop of %d\n" , count);
			})
		}
	}()
	go func() {
		defer wg.Done()
		for count<100{
			oneCond.DoFunc(false,func() {
					count++
					fmt.Printf("go 协程 1,loop of %d\n" , count);
			})
		}
	}()
	wg.Wait()
}

func TestCond2(t *testing.T){
	var wg sync.WaitGroup
	wg.Add(2)
	twoCond:=NewTwoCondition()
	count:=0

	go func() {
		defer wg.Done()
		for count<100{
			twoCond.DoNextFunc(func() {
				count++
				fmt.Printf("go 协程 2,loop of %d\n" , count);
			})
		}
	}()
	go func() {
		defer wg.Done()
		for count<100{
			twoCond.DoFrontFunc(func() {
				count++
				fmt.Printf("go 协程 1,loop of %d\n" , count);
			})
		}
	}()

	wg.Wait()

}



/*
N个线程，按顺序交替执行打印,要求结果如下：
*/

const (
	N      = 2
)


func TestConds(t *testing.T) {

}

func BenchmarkSimpleWorkerPool(b *testing.B) {
	//for i:=0;i< b.N;i++{

	//}
}
