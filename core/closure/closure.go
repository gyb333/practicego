package closure

import "fmt"

/*
在程序编译完后，函数会以计算机指令⽅式存储在代码区，在定义函数时指定的形参，在未出现函数调用时，它们并不占内存中的存储单元，
函数的参数以及局部变量会在调用时加载到内存中。
栈总是向下增长的。压栈的操作使得栈顶的地址减小，弹出操作使得栈顶地址增⼤。
函数调用过程所需要的信息：
	函数的返回地址；
	函数的参数；
	保存的上下⽂：包括在函数调用前后需要保持不变的寄存器。

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
type it=int		//it 声明的类型可以直接计算

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


