package queue

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

// 队列接口
type Queue interface {
	IsRunning() bool
	SetRunning(isRunning bool)

	Put(e interface{})

	Swap(msgBuf *MessageBuffer) *MessageBuffer

	IsEmpty() bool
}


// 队列结构体
type CappedBlockingQueue struct {
	msgQueue *MessageBuffer
	cap      int
	running  bool

	mu    *sync.Mutex
	full  *sync.Cond
	empty *sync.Cond
}

func NewCappedBlockingQueue(cap int) Queue {
	queue := &CappedBlockingQueue{
		msgQueue: NewMessageBuffer(),
		cap:      cap,
		running:  true,
		mu: new(sync.Mutex),
	}
	// 两个条件变量公用一个锁
	queue.full = sync.NewCond(queue.mu)
	queue.empty = sync.NewCond(queue.mu)
	return queue
}

func (q CappedBlockingQueue) IsRunning() bool {
	return q.running
}

func (q *CappedBlockingQueue) SetRunning(isRunning bool) {
	q.running = isRunning
}

func (q *CappedBlockingQueue) Put(e interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()
	// 循环检查不可替换为if
	for q.msgQueue.Size() >= q.cap && q.running {
		q.full.Wait() // 阻塞goroutine并释放锁
	}
	if !q.running {
		return
	}
	q.msgQueue.Add(e)
	q.empty.Signal() // 唤醒单个goroutine消费
}

func (q *CappedBlockingQueue) Swap(msgBuf *MessageBuffer) *MessageBuffer {
	q.mu.Lock()
	defer q.mu.Unlock()
	for q.msgQueue.IsEmpty() && q.running {
		q.empty.Wait() // 阻塞goroutine并释放锁
	}
	toReturn := q.msgQueue
	q.msgQueue = msgBuf
	q.full.Broadcast() // 唤醒所有goroutine生产
	return toReturn
}

func (q CappedBlockingQueue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.msgQueue.IsEmpty()
}


//测试代码
func main() {
	queue := NewCappedBlockingQueue(20)

	for i:=0; i<10;i++  {
		go func(i int) {
			for {
				for j:=0; j<10; j++  {
					queue.Put("produce by goroutine -> " + strconv.Itoa(i))
				}
			}
		}(i)
	}

	for k:=0; k<5; k++ {
		go func(i int) {
			for {
				newBuf := NewMessageBuffer()
				oldBuf := queue.Swap(newBuf)
				for v := oldBuf.Front(); v!=nil; v = v.Next()  {
					fmt.Println("consume by ", i, " ", v.Value)
				}
			}
		}(k)
	}

	time.Sleep(2*time.Minute)
}