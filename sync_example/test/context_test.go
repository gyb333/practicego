package test_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func iscanceled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func dowork() {
	fmt.Println("do work")
	time.Sleep(time.Second)
}

func TestCancel(t *testing.T) {
	wg:=sync.WaitGroup{}
	N:=5
	wg.Add(N)
	ctx, cancelfunc := context.WithCancel(context.Background())
	for i := 0; i < N; i++ {
		go func(i int, ctx context.Context) {
			defer wg.Done()
			for {
				if iscanceled(ctx) {
					break
				}
				dowork()
			}
			fmt.Println(i, "canceled.")
		}(i, ctx)
	}
	time.Sleep(1*time.Second)
	cancelfunc()
	wg.Wait()
}
