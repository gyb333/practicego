package memory_mode

import (
	"fmt"
	"unsafe"
)
type FUNC func(int)int		//1个机器字


func DataType()  {
	var a struct{}			//内存空间为0，在栈创建的变量指向一个固定地址
	fmt.Printf("%p,%T,%v,%d\n",&a,a,a,unsafe.Sizeof(a))	//0x685358,struct {},{},0
	b:= struct {}{}
	fmt.Printf("%p,%T,%v,%d\n",&b,b,b,unsafe.Sizeof(b))	//0x685358,struct {},{},0

	var i interface{}
	fmt.Printf("%p,%T,%v,%d\n",&i,i,i,unsafe.Sizeof(i))	//0xc0000384e0,<nil>,<nil>,16
	var ie interface{}
	fmt.Printf("%p,%T,%v,%d\n",&ie,ie,ie,unsafe.Sizeof(ie))	//0xc0000384f0,<nil>,<nil>,16
	var f FUNC
	fmt.Printf("%p,%T,%v,%d\n",&f,f,f,unsafe.Sizeof(f))
	
}
var vx struct {
	i int8
	b bool
	f float32	//4 结构体内存为最大字节整数倍
	r rune
}
type x struct {
	i int8
	b bool
	f float32
	r rune
}

type y struct {
	x				//x的12  根据下面一个类型为int32所以为4字节 不需要扩充，如果是int64 则扩充到16
	a [5]int32		//5*4 20 	偏移位一定是下一类型不超过操作系统的机器字(32位 4 64位 8)的整数倍，[4],[5],都是一样的内存分配; [3] 则3*4=12 加上上面12 正好是8的倍数，不用对齐
	s string		//2*8 16
	sc []int		//3*8 24
	m map[string]interface{}	//1个机器字*8 8
	p *int				//1个机器字*8 8
	c chan struct{}		//1个机器字 8
}

func StructAlignment()  {
	fmt.Printf("%p,%v,%v\n",&vx,unsafe.Sizeof(vx), unsafe.Alignof(vx))//全局变量0x6853b8,12,4
	x:=x{}
	fmt.Printf("%p,%v,%v\n",&x,unsafe.Sizeof(x), unsafe.Alignof(x))	//0xc00000a3b0,12,4
	y:=y{}
	fmt.Printf("%p,%v,%v\n",&y,unsafe.Sizeof(y), unsafe.Alignof(y))	//0xc000044180,96,8
	fmt.Println("-----------------------------")
	fmt.Printf("%p,%v,%v\n",&y.i,unsafe.Sizeof(y.i),  unsafe.Offsetof(y.i))	//0xc000044180,1,0
	fmt.Printf("%p,%v,%v\n",&y.a,unsafe.Sizeof(y.a),  unsafe.Offsetof(y.a))	//0xc000044190,16,16
	fmt.Printf("%p,%v,%v\n",&y.s,unsafe.Sizeof(y.s),  unsafe.Offsetof(y.s))	//0xc0000441a0,16,32	32=上一行的16+16
	fmt.Printf("%p,%v,%v\n",&y.sc,unsafe.Sizeof(y.sc),  unsafe.Offsetof(y.sc))	//0xc0000441b0,24,48 48=上一行的16+32
	fmt.Printf("%p,%v,%v\n",&y.m,unsafe.Sizeof(y.m),  unsafe.Offsetof(y.m))	//0xc0000441c8,8,72
	fmt.Printf("%p,%v,%v\n",&y.p,unsafe.Sizeof(y.p),  unsafe.Offsetof(y.p))	//0xc0000441d0,8,80
	fmt.Printf("%p,%v,%v\n",&y.c,unsafe.Sizeof(y.c),  unsafe.Offsetof(y.c))	//0xc0000441d8,8,88
}
