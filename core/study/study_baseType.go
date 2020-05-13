package basic

import (
	"fmt"
	"unsafe"
)

//基本数据类型：byte int uint int8 uint8 int16 uint16 int32 uint32 int64 uint64 uintptr
//   float32 float64 complex64 complex128
//   bool rune string error
func BaseDataType()  {
	var a byte=33
	fmt.Printf("%p,%d,%v,%d,%#x,%T\n",&a,unsafe.Sizeof(a),a,a,a,a)
	var i =33
	fmt.Printf("%p,%d,%v,%d,%#x,%T\n",&i,unsafe.Sizeof(i),i,i,i,i)
	var i32 int32 =33
	fmt.Printf("%p,%d,%v,%d,%#x,%T\n",&i32,unsafe.Sizeof(i32),i32,i32,i32,i32)

	var f float64 = 33
	fmt.Printf("%p,%d,%v,%f,%T\n",&f,unsafe.Sizeof(f),f,f,f)

	var b =true
	fmt.Printf("%p,%d,%v,%t,%T\n",&b,unsafe.Sizeof(b),b,b,b)



}