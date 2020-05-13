package dataType

import (
	"fmt"
	"reflect"
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
unsafe.Alignof 函数返回对应参数的类型需要对齐的倍数. 和 Sizeof 类似, Alignof 也是返回一个常量表达式, 对应一个常量.
通常情况下布尔和数字类型需要对齐到它们本身的大小(最多8个字节), 其它的类型对齐到机器字大小.

unsafe.Offsetof 函数的参数必须是一个字段 x.f, 然后返回 f 字段相对于 x 起始地址的偏移量, 包括可能的空洞.
**/

func BaseStruct()  {
	var b bool
	var i8 int8
	var u8 uint8
	var i16 int16
	var u16 uint16
	var i32 int32
	var u32 uint32
	var i int
	var u uint
	var f32 float32
	var f64 float64
	var ptr uintptr
	var p *int
	var str string
	var array [7]byte
	var arr [7]rune
	var slice []int
	var m map[string]int
	var c  chan struct{}
	f :=func (){}		//匿名函数
	var s struct{}		//匿名结构体
	var in interface{}	//匿名接口

	fmt.Println("var b bool			//",reflect.TypeOf(b),Sizeof(b),Alignof(b))					//bool 1 1
	fmt.Println("var i8 int8 		//",reflect.TypeOf(i8),Sizeof(i8),Alignof(i8))				//int8 1 1
	fmt.Println("var u8 uint8 		//",reflect.TypeOf(u8),Sizeof(u8),Alignof(u8))				//uint8 1 1
	fmt.Println("var i16 int16		//",reflect.TypeOf(i16),Sizeof(i16),Alignof(i16))			//int16 2 2
	fmt.Println("var u16 uint16		//",reflect.TypeOf(u16),Sizeof(u16),Alignof(u16))			//uint16 2 2
	fmt.Println("var i32 int32		//",reflect.TypeOf(i32),Sizeof(i32),Alignof(i32))			//int32 4 4
	fmt.Println("var u32 uint32		//",reflect.TypeOf(u32),Sizeof(u32),Alignof(u32))			//uint32 4 4
	fmt.Println("var i int			//",reflect.TypeOf(i),Sizeof(i),Alignof(i))					//int 8 8
	fmt.Println("var u uint			//",reflect.TypeOf(u),Sizeof(u),Alignof(u))					//uint 8 8
	fmt.Println("var f32 float32		//",reflect.TypeOf(f32),Sizeof(f32),Alignof(f32))			//float32 4 4
	fmt.Println("var f64 float64		//",reflect.TypeOf(f64),Sizeof(f64),Alignof(f64))			//float64 8 8

	fmt.Println("var ptr uintptr		//",reflect.TypeOf(ptr),Sizeof(ptr),Alignof(ptr))			//uintptr 8 8
	fmt.Println("var p *int			//",reflect.TypeOf(p),Sizeof(p),Alignof(p))					//*int 8 8

	fmt.Println("var str string		//",reflect.TypeOf(str),Sizeof(str),Alignof(str))			//string 16 8
	fmt.Println("var array [7]byte	//",reflect.TypeOf(array),Sizeof(array),Alignof(array))		//[8]uint8 8 1
	fmt.Println("var arr [7]rune	//",reflect.TypeOf(arr),Sizeof(arr),Alignof(arr))		//[7]uint32 28 4

	fmt.Println("var slice []int		//",reflect.TypeOf(slice),Sizeof(slice),Alignof(slice))		//[]int 24 8
	fmt.Println("var m map[string]int	//",reflect.TypeOf(m),Sizeof(m),Alignof(m))					//map[string]int 8 8
	fmt.Println("var c  chan struct{}	//",reflect.TypeOf(c),Sizeof(c),Alignof(c))					//chan struct {} 8 8

	fmt.Println("f :=func (){}			//",reflect.TypeOf(f),Sizeof(f),Alignof(f))					//func() 8 8
	fmt.Println("var s struct{}			//",reflect.TypeOf(s),Sizeof(s),Alignof(s))					//struct {} 0 1
	fmt.Println("var in interface{}		//",reflect.TypeOf(in),Sizeof(in),Alignof(in))				//<nil> 16 8

}


//内存对齐：操作系统位数的整数倍，即8字节整数倍
var x struct {	//匿名结构体类型
	a bool			//1个字节
	b int16			//2个字节	a和b 合并成一个机器字
	c []int			//3个机器字 24字节
	string		//匿名字段默认把类型设置为名称	2个机器字 16字节

}

func StructOffset() {
	fmt.Printf("%p,%v,%v\n",&x,Sizeof(x), Alignof(x))	//(1+3+2)*8=48,1+2+3+2=8
	fmt.Printf("%p,%v,%v,%v\n",&x.a,Sizeof(x.a), Alignof(x.a), Offsetof(x.a))	//1		1	0
	fmt.Printf("%p,%v,%v,%v\n",&x.b,Sizeof(x.b), Alignof(x.b), Offsetof(x.b))	//2		2	2
	fmt.Printf("%p,%v,%v,%v\n",&x.c,Sizeof(x.c), Alignof(x.c), Offsetof(x.c))	//24	8	8
	fmt.Printf("%p,%v,%v,%v\n",&x.string,Sizeof(x.string), Alignof(x.string), Offsetof(x.string))	//16	8	32

	x.string="string"
	fmt.Println(x.string) // "42"
	// 和 pb := &x.b 等价
	pb := (*int16)(Pointer(uintptr(Pointer(&x)) + Offsetof(x.b)))
	*pb = 42
	fmt.Println(pb,x.b) // "42"
	pb=&x.b
	*pb=11
	fmt.Println(pb,*pb,x.b)
	fmt.Printf("%b\n",pb)

}

