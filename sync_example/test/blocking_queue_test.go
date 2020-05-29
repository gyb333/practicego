package test_test

import (
	"fmt"
	"sync"
	"testing"
)


type SimpleBlockingQueue struct {
	data chan interface{}
	length int
}

func NewSimpleBlockingQueue(length int)  *SimpleBlockingQueue{
	return &SimpleBlockingQueue{
		data: make(chan interface{},length),
		length: length,
	}
}

func (q *SimpleBlockingQueue) Put(data interface{})  {
	q.data<-data
}
func (q *SimpleBlockingQueue) Get()  interface{}{
	return <-q.data
}


func TestBlockingQueue(t *testing.T) {
	wg :=sync.WaitGroup{}
	wg.Add(2)
	q:=NewSimpleBlockingQueue(10)
	count:=0
	go func() {
		defer wg.Done()
		for count<=100{
			count++
			q.Put(count)
		}
	}()
	go func() {
		defer wg.Done()
		for count<=100{
			fmt.Println(q.Get())
		}
	}()
	wg.Wait()
}
