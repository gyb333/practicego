package sync_example_test

import (
	"fmt"
	"sync"
	"testing"
)

/*
请编写2个线程，线程1顺序输出1，3，5，……, 99 等奇数，每个数 一 。
线程2顺序输出2，4，6……100等偶数，每个数 一 。
最终的结果要求是输出为 自然顺序：1, 2, 3, 4, ……99, 100。
 */

func Test2(t *testing.T)  {
	var wg sync.WaitGroup
	wg.Add(2)
	c1 :=make(chan int)
	c2 :=make(chan int)

	go func() {
		defer wg.Done()
		for j:=1;j<=50;j++{
			fmt.Printf("协程1,输出：%d\n",<-c1)
			c2<-j+j
			}

	}()
	go func() {
		defer wg.Done()
		for j:=1;j<=50;j++{
			c1<-j+j-1
			fmt.Printf("协程2,输出：%d\n",<-c2)
		}
	}()
	wg.Wait()
}
/*
N个线程，按顺序交替执行打印,要求结果如下：
 */

func TestN(t *testing.T)  {
	var wg sync.WaitGroup
	N:=5
	wg.Add(N)
	var cs [5]chan int
	for c :=range cs{
		c = make(chan int)
	}


	go func() {
		defer wg.Done()
		for j:=1;j<=50;j++{
			fmt.Printf("协程1,输出：%d\n",<-c1)
			c2<-j+j
		}

	}()
	go func() {
		defer wg.Done()
		for j:=1;j<=50;j++{
			c1<-j+j-1
			fmt.Printf("协程2,输出：%d\n",<-c2)
		}
	}()
	wg.Wait()
}
