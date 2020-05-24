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

func Test(t *testing.T) {
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

func TestCond(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	cond := sync.Cond{L: new(sync.Mutex)}
	//让协程先执行等待
	go func() {
		defer wg.Done()
		cond.L.Lock()
		for i := 1; i <= 50; i++ {
			cond.Wait()
			println("协程g2:", i+i) //执行步骤4
			if i < 50 {
				cond.Signal()
			}
		}
		cond.L.Unlock()
	}()

	go func() {
		defer wg.Done()
		cond.L.Lock()
		for i := 1; i <= 50; i++ {
			cond.Signal()
			println("协程g1:", i+i-1) // 执行步骤1， 执行步骤5
			if i < 50 {
				cond.Wait()
			}
		}
		cond.L.Unlock()
	}()
	wg.Wait()
}

func TestCond1(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	cond := sync.Cond{L: new(sync.Mutex)}
	for i := 1; i <= 2; i++ {
		//让协程先执行等待
		go func(j int) {
			defer wg.Done()
			cond.L.Lock()
			for i := 1; i <= 50; i++ {
				if j == 1 {
					cond.Wait()
					println("协程g2:", i+i) //执行步骤4
					cond.Signal()
				} else {
					cond.Signal()
					println("协程g1:", i+i-1) //执行步骤4
					cond.Wait()
				}
			}
			cond.L.Unlock()
		}(i)
	}

	wg.Wait()
}

/*
N个线程，按顺序交替执行打印,要求结果如下：
*/

const (
	N      = 3
	COUNTS = 100
)

func cond2() {
	var wg sync.WaitGroup
	wg.Add(N)
	ca :=sync.Cond{L: new(sync.Mutex)}
	cb :=sync.Cond{L: new(sync.Mutex)}

	go func() {
		defer wg.Done()
		ca.L.Lock()
		for i := 1; i <= 50; i++ {
			ca.Wait()
			println("协程g2:", i+i) //执行步骤4
			cb.Signal()
		}
		ca.L.Unlock()
	}()

	go func() {
		defer wg.Done()
		cb.L.Lock()
		for i := 1; i <= 50; i++ {
			ca.Signal()
			println("协程g1:", i+i-1) // 执行步骤1， 执行步骤5
			cb.Wait()
		}
		cb.L.Unlock()
	}()


	wg.Wait()

}


func cond3() {
	var wg sync.WaitGroup
	wg.Add(N)
	ca :=sync.Cond{L: new(sync.Mutex)}
	cb :=sync.Cond{L: new(sync.Mutex)}
	cc :=sync.Cond{L: new(sync.Mutex)}
	go func() {
		defer wg.Done()
		cc.L.Lock()
		for i := 1; i <= 50; i++ {
			cc.Wait()
			println("协程g3:", i*3) //执行步骤4
			ca.Signal()
		}
		cc.L.Unlock()
	}()

	go func() {
		defer wg.Done()
		cb.L.Lock()
		for i := 1; i <= 50; i++ {
			cb.Wait()
			println("协程g2:", i*3-1) //执行步骤4
			cc.Signal()
		}
		cb.L.Unlock()
	}()

	go func() {
		defer wg.Done()
		ca.L.Lock()
		for i := 1; i <= 50; i++ {
			cb.Signal()
			println("协程g1:", i*3-2) // 执行步骤1， 执行步骤5
			ca.Wait()
		}
		ca.L.Unlock()
	}()


	wg.Wait()

}

func conds() {
	var wg sync.WaitGroup
	wg.Add(N)
	var cs [N]sync.Cond
	for i := range cs {
		cs[i] = sync.Cond{L: new(sync.Mutex)}
	}
	for j := range cs {
		go func() {
			defer wg.Done()
			cs[j].L.Lock()
			for i := 1; i <= 50; i++ {
				cs[j].Wait()
				fmt.Printf("协程g%d:%d\n",j, i) //执行步骤4
				cs[(j+1)%N].Signal()
			}
			cs[j].L.Unlock()
		}()
	}
	wg.Wait()

}
func TestConds(t *testing.T) {
	cond3()
}

func BenchmarkSimpleWorkerPool(b *testing.B) {
	//for i:=0;i< b.N;i++{
	conds()
	//}
}
