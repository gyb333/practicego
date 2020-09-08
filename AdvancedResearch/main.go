package main

import "fmt"

func main()  {
	a := make(chan int, 1)
	b := make(chan int, 1)
	c := make(chan int, 1)

	//a>b>c
	//使用select中的default特性实现优先级队列。
	for {
		select {
		case <-a:
			fmt.Println("recv from a")
		default:
			select {
			case <-b:
				fmt.Println("recv from b")
			default:
				select {
				case <-c:
					fmt.Println("recv from c")
				}
			}
		}

	}


}
