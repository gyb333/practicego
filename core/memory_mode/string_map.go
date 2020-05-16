package memory_mode

import (
	"fmt"
	"reflect"
	"unsafe"
)

type Student struct {
	name string
	age int8
}

func StringMemory()  {
	str :="雇佣兵333"
	p := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&str)).Data)
	fmt.Printf("%p,%p\n",&str,p)

	bs := []byte(str)
	p =unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&bs)).Data)
	fmt.Printf("%p,%p\n",bs,p)

	s:=string(bs)
	p = unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&s)).Data)
	fmt.Printf("%p,%p\n",&s,p)
}

func MapMemory()  {
	m :=make(map[int]Student)
	s :=Student{name: "GYB",age: 33}
	m[1001]=s		//深拷贝
	s.age++
	fmt.Printf("%v,%v\n",s, m[1001])
}