package closure

import "fmt"

/*
因为这个匿名函数捕获了外部函数的局部变量v，这种函数我们一般叫闭包
defer 后进先出
*/

func Closure(){
	f, g := fa(0)
	s, k := fa(0)
	fmt.Println(f(1), g(2))		//0+1=1		1-2=-1
	fmt.Println(f(3), g(2))		//-1+3=2	2-2=0

	fmt.Println(s(1), k(2))		//0+1=1		1-2=-1
}
type FUNC func(int)int		//没有=不能直接计算

func fa(base int) (FUNC, func(int) int) {
	fmt.Println(&base, base)
	add := func(i int) int {
		base += i
		fmt.Println("add:",&base, base)
		return base
	}
	sub := func(i int) int {
		base -= i
		fmt.Println("sub:",&base, base)
		return base
	}
	return add, sub
}


