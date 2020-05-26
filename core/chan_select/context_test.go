package chan_select_test

import (
	"context"
	"fmt"
	"testing"
	"time"
)
type MyContext struct {
	// 这里的 Context 是我 copy 出来的，所以前面不用加 context.
	context.Context
}

func Test2(t *testing.T) {
	childCancel := true
	parentCtx, parentFunc := context.WithCancel(context.Background())
	mctx := MyContext{parentCtx}
	childCtx, childFun := context.WithCancel(mctx)
	if childCancel {
		childFun()
	} else {
		parentFunc()
	}
	go func() {
		tc:= time.NewTicker(10*time.Nanosecond)
	L:
		for{
			select {
			case <-tc.C:
			case <-parentCtx.Done():
				fmt.Println("Parent done")
			case <-childCtx.Done():
				fmt.Println("Child done")
				break L
			}
		}
	}()

	fmt.Println(parentCtx)
	fmt.Println(mctx)
	fmt.Println(childCtx)
	// 防止主协程退出太快，子协程来不及打印
	time.Sleep(1 * time.Second)
}
func process(ctx context.Context) {
	traceId, ok := ctx.Value("traceId").(string)
	if ok {
		fmt.Printf("process over. trace_id=%s\n", traceId)
	} else {
		fmt.Printf("process over. no trace_id\n")
	}
}
func TestContext(t *testing.T) {
	ctx := context.Background()
	process(ctx)
	ctx = context.WithValue(ctx, "traceId", "qcrao-2019")
	process(ctx)
}


func gen(ctx context.Context) <-chan int {
	ch := make(chan int)
	go func() {
		var n int
		for {
			select {
			case <-ctx.Done():
				return
			case ch <- n:
				n++
				time.Sleep(time.Second)
			}
		}
	}()
	return ch
}

func TestCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 避免其他地方忘记 cancel，且重复调用不影响
	for n := range gen(ctx) {
		fmt.Println(n)
		if n == 5 {
			cancel()
			break
		}
	}
}