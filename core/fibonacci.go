package main

import (
	"math/big"
)

func Fibonacci(n int) *big.Int {
	x, y := big.NewInt(0), big.NewInt(1)
	for i := 0; i < n; i++ {
		x, y = y, x.Add(x, y)
	}
	return x
}

func fibonacciChan(n int) (result int) {
	channel := make(chan int)
	quit := make(chan bool)
	go func() {
		for i := 0; i < n; i++ {
			result = <-channel
		}
		quit <- true
	}()
	func(channel chan int, quit chan bool) {
		x, y := 0, 1
		for {
			select {
			case channel <- y:
				x, y = y, x+y
			case <-quit:
				return
			}
		}
	}(channel, quit)
	return
}

func FibonacciChanBig(n int) (result *big.Int) {
	result = big.NewInt(0)
	if n == 0 {
		return
	}
	channel := make(chan *big.Int)
	quit := make(chan bool)
	go func() {
		for i := 0; i < n; i++ {
			result = <-channel
		}
		quit <- true
	}()
	func(channel chan *big.Int, quit chan bool) {
		x, y := big.NewInt(0), big.NewInt(1)
		for {
			select {
			case channel <- y:
				x, y = y, x.Add(x, y)
			case <-quit:
				return
			}
		}
	}(channel, quit)
	return
}

func FibonacciClosure(n int) (r *big.Int) {
	if n < 1 {
		return big.NewInt(0)
	}
	f := func() func() *big.Int {
		x, y := big.NewInt(0), big.NewInt(1)
		return func() *big.Int {
			x, y = y, x.Add(x, y)
			return x
		}
	}()
	for i := 0; i < n; i++ {
		r = f()
	}
	return
}

func fibonacciBig(n int) (r *big.Int) {
	if n < 2 {
		return big.NewInt(int64(n))
	}
	f := func() func() *big.Int {
		v, s := big.NewInt(0), big.NewInt(1)
		return func() *big.Int {
			var tmp big.Int
			tmp.Set(s)
			s.Add(s, v)
			v = &tmp
			return s
		}
	}()

	for i := 1; i < n; i++ {
		r = f()
	}
	return

}

func FibonacciRecursion(n int) (res int) {
	if n <= 1 {
		res = n
	} else {
		res = FibonacciRecursion(n-2) + FibonacciRecursion(n-1)
	}
	return
}
