package locks

import (
	"context"
	"fmt"
	"time"
)

func ExampleUnboundedExecutor_Go() {
	executor := NewUnboundedExecutor()
	executor.Go(func(ctx context.Context) {
		fmt.Println("abc")
	})
	time.Sleep(time.Second)
	// output: abc
}

func ExampleUnboundedExecutor_StopAndWaitForever() {
	executor := NewUnboundedExecutor()
	executor.Go(func(ctx context.Context) {
		everyMillisecond := time.NewTicker(time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("goroutine exited")
				return
			case <-everyMillisecond.C:
				// do something
			}
		}
	})
	time.Sleep(time.Second)
	executor.StopAndWaitForever()
	fmt.Println("executor stopped")
	// output:
	// goroutine exited
	// executor stopped
}

func ExampleUnboundedExecutor_Go_panic() {
	HandlePanic = func(recovered interface{}, funcName string) {
		fmt.Println(funcName)
	}
	executor := NewUnboundedExecutor()
	executor.Go(willPanic)
	time.Sleep(time.Second)
	// output:
	// github.com/modern-go/concurrent_test.willPanic
}

func willPanic(ctx context.Context) {
	panic("!!!")
}
