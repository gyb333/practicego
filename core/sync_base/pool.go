package sync_base

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
)

type Pool struct {
	m sync.Mutex
	res chan io.Closer
	//创建资源的方法，由用户程序自己生成传入
	factory func() (io.Closer, error)
	closed bool
	//资源池获取资源超时时间
	timeout <-chan time.Time
}
//资源池关闭标志
var ErrPoolClosed =  errors.New("资源池已经关闭")
//超时标志
var ErrTimeout = errors.New("获取资源超时")
//新建资源池
func New(fn func() (io.Closer, error), size int) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("新建资源池大小太小")
	}
	//新建资源池
	p := Pool{
		factory: fn,
		res: make(chan io.Closer, size),
	}
	//向资源池循环添加资源，直到池满
	for count := 1; count <= cap(p.res); count++ {
		r, err := fn()
		if err != nil {
			fmt.Println("添加资源失败，创建资源方法返回nil")
			break
		}
		fmt.Println("资源加入资源池")
		p.res <- r
	}
	fmt.Println("资源池已满，返回资源池")
	return &p, nil
}
//获取资源
func (p *Pool) Acquire(d time.Duration) (io.Closer, error) {
	//设置d时间后超时
	p.timeout = time.After(d)
	select {
	case r, ok := <-p.res:
		fmt.Println("获取", "共享资源")
		if !ok {
			return nil, ErrPoolClosed
		}
		return r, nil
	case <-p.timeout:
		return nil, ErrTimeout
	}
}
//放回资源池
func (p *Pool) Release(r io.Closer) {
	//上互斥锁，和Close方法对应，不同时操作
	p.m.Lock()
	defer p.m.Unlock()
	if p.closed {
		r.Close()
		return
	}
	//资源放回队列
	select {
	case p.res <- r:
		fmt.Println("资源放回队列")
	default:
		fmt.Println("资源队列已满，释放资源")
		r.Close()
	}
}
//关闭资源池
func (p *Pool) Close() {
	//互斥锁，保证同步，和Release方法相关，用同一把锁
	p.m.Lock()
	defer p.m.Unlock()
	if p.closed {
		return
	}
	p.closed = true
	//清空通道资源之前，将通道关闭，否则引起死锁
	close(p.res)
	for r := range p.res {
		r.Close()
	}
}
