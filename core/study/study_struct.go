package basic

import (
	"fmt"
	. "unsafe"
)

/**
类型															大小
bool														1个字节
intN, uintN, floatN, complexN							N/8个字节(例如float64是8个字节)
int, uint, uintptr										1个机器字
*T														1个机器字
string													2个机器字(data,len)
[]T														3个机器字(data,len,cap)
map														1个机器字
func													1个机器字
chan													1个机器字
interface												2个机器字(type,value)
unsafe.Alignof 函数返回对应参数的类型需要对齐的倍数. 和 Sizeof 类似, Alignof 也是返回一个常量表达式, 对应一个常量. 通常情况下布尔和数字类型需要对齐到它们本身的大小(最多8个字节), 其它的类型对齐到机器字大小.

unsafe.Offsetof 函数的参数必须是一个字段 x.f, 然后返回 f 字段相对于 x 起始地址的偏移量, 包括可能的空洞.
**/



var x struct {	//匿名结构体类型
	a bool
	b int16
	c []int
	string		//匿名字段默认把类型设置为名称
}

func structOffset() {
	fmt.Println(Sizeof(x), Alignof(x))
	fmt.Println(Sizeof(x.a), Alignof(x.a), Offsetof(x.a))
	fmt.Println(Sizeof(x.b), Alignof(x.b), Offsetof(x.b))
	fmt.Println(Sizeof(x.c), Alignof(x.c), Offsetof(x.c))
	// 和 pb := &x.b 等价
	pb := (*int16)(Pointer(uintptr(Pointer(&x)) + Offsetof(x.b)))
	*pb = 42
	fmt.Println(x.b,x.string) // "42"
}


