package defer_recover

import "fmt"

/*
defer延迟调用:
	释放占用的资源
	捕捉处理异常
	输出日志
 */
/*
因为这个匿名函数捕获了外部函数的局部变量v，这种函数我们一般叫闭包
defer 后进先出
*/
func DeferFunc() {
	for i := 0; i < 3; i++ {
		defer func() { fmt.Println("first:",i) }() 		//3		3		3
	}
	for i := 0; i < 3; i++ {
		i := i // 定义一个循环体内局部变量i
		defer func() { fmt.Println("second:",i) }()		//2		1		0
	}
	for i := 0; i < 3; i++ {
		// 通过函数传入i defer 语句会马上对调用参数求值
		defer func(i int) { fmt.Println("threed:",i) }(i)		//2		1		0
	}

}



/*
recover错误拦截：
	编辑时异常
	编译时异常
	运行时异常
*/
func PanicRecover() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover ", err)	//recover  first defer panic
		}
	}()	//延时调用  函数退出或panic触发
	defer func() {
		panic("first defer panic")
	}()
	defer func() {
		panic("second defer panic ")
	}()
	panic("main body panic ")
	panic("main1 body panic ")	//不会执行
}