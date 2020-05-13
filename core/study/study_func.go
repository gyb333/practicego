package basic

import "fmt"

/*
值传递：操作的是数值本身
值类型：存储数值本身,所有的基本类型 数组 结构体
引用传递：操作的是数值的地址
	修改变量实践是修改变量地址处的内存

引用类型:存储数值的地址 slice map chan func pointer interface
 */


func FuncMain()  {

	//args :=[]interface{}{1,2,3,"4",nil}
	//variableFunc(args...)
	//closure()
	//deferFunc()
	panicRecover()
}


func variableFunc(args ...interface{}){
	if args !=nil{
		for i,v :=range args{
			fmt.Println(i,v)
		}
	}
}


func closure(){
	f, g := fa(0)
	s, k := fa(0)
	fmt.Println(f(1), g(2))		//0+1=1		1-2=-1
	fmt.Println(f(3), g(2))		//-1+3=2	2-2=0
	fmt.Println(s(1), k(2))		//0+1=1		1-2=-1
}

func fa(base int) (func(int) int, func(int) int) {
	fmt.Println(&base, base)
	add := func(i int) int {
		base += i
		fmt.Println(&base, base)
		return base
	}
	sub := func(i int) int {
		base -= i
		fmt.Println(&base, base)
		return base
	}
	return add, sub
}

/*
因为这个匿名函数捕获了外部函数的局部变量v，这种函数我们一般叫闭包
defer 后进先出
 */
func deferFunc() {
	for i := 0; i < 3; i++ {
		defer func() { fmt.Println(i) }() 		//3		3		3
	}
	for i := 0; i < 3; i++ {
		i := i // 定义一个循环体内局部变量i
		defer func() { fmt.Println(i) }()		//2		1		0
	}
	for i := 0; i < 3; i++ {
		// 通过函数传入i defer 语句会马上对调用参数求值
		defer func(i int) { fmt.Println(i) }(i)		//2		1		0
	}

}

func panicRecover() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover ", err)
		}
	}()
	defer func() {
		panic("first defer panic")
	}()
	defer func() {
		panic("second defer panic ")
	}()
	panic("main body panic ")
}