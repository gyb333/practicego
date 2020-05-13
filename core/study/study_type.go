package basic

import (
	"fmt"
	"unsafe"
		"bytes"
	"encoding/binary"
)


/**
值类型 深拷贝：基本数据类型  数组 结构体
引用类型 浅拷贝: 切片 字典 通道  函数 指针类型
 */

func init() {
	fmt.Println("basic package init")
}

const (
	Monday    = iota //0
	Tuesday          //默认自动加1 1
	Wednesday        //2
	Thursday
	Friday
	Saturday
	Sunday
)
const (
	_          = iota
	KB float64 = 1 << (iota * 10)
	MB
	GB
	TB
	PB
)




const ptrIntSize = unsafe.Sizeof((*int)(nil))
const ptrIntAlign = unsafe.Alignof((*int)(nil))

func unsafeData()  {
	fmt.Println(ptrIntAlign,ptrIntSize)
}

type typeAlias = int	//类型别名
type SliceInt []int		//定义一个新的类型名称

func (s SliceInt) Sum() int {
	sum := 0
	for _, i := range s {
		sum += i
	}
	return sum
}
func SliceInt_Sum(s SliceInt) int {
	sum := 0
	for _, i := range s {
		sum += i
	}
	return sum
}


func BasicType()  {
	unsafeData()
	var s SliceInt = []int{1, 2, 3, 4}
	println(s.Sum())
	println(SliceInt_Sum(s))
}









func toBytes() []byte {
	bb :=bytes.NewBuffer(nil)
	binary.Write(bb, binary.BigEndian, 'h')
	bs := bb.Bytes()
	fmt.Printf("%#X,%d,%d\n",bs,len(bs),cap(bs))
	return bs
}




