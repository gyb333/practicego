package test_test

import (
	. ".."
	"fmt"
	"sync"
	"testing"
)

/*
N个线程，按顺序交替执行打印,要求结果如下：
*/



const (
	  COUNTS  =1000
	  GOS     =10
)
func TestConditions(t *testing.T) {
	wg:=sync.WaitGroup{}
	wg.Add(GOS)
	c:=NewNConditions(GOS)
	count:=0
	for i:=0;i<c.Length;i++{
		go func(i int) {
			defer wg.Done()
			for count<=COUNTS{
				c.DoFunc(i, func() {
					count++
					if count<=COUNTS{
						fmt.Printf("协程%2d,执行数据%4d\n",i+1,count)
					}
				})
			}
		}(i)
	}
	wg.Wait()
}