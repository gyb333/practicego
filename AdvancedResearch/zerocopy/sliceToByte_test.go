package zerocopy

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unsafe"
)

func Test_NoCopy(t *testing.T)  {
	str:="Hello 世界！" //字符串常量 strings.Repeat("abc",3) 底层数据可以修改
	sh := (*reflect.StringHeader)(unsafe.Pointer(&str))
	p := unsafe.Pointer(sh.Data)
	fmt.Println(&str, str,  sh, p)
	//as :=(*[15]byte)(p)
	//bas :=as[:]
	//copy(bas,[]byte(strings.Repeat("h",2)))	//底层切片也不允许被修改
	//for i,v :=range as{
	//	as[i]=v+1	//底层数组不允许被修改
	//}
	fmt.Println(str)


	//bs:=stringToBytes(str)
	bs:=stringToBytes2(str)
	//bs:=stringToByteSlice(str)
	bsh :=(*reflect.SliceHeader)(unsafe.Pointer(&bs))
	p =unsafe.Pointer(bsh.Data)
	fmt.Println(unsafe.Pointer(&bs),bs,*bsh,p)

	bss :=bs[0:]
	//copy(bss,[]byte(strings.Repeat("h",17)))//长度没有变化,底层数组不允许被遍历修改
	bss=append(bss,[]byte(strings.Repeat("h",3))...)	//添加的数据超过底层数组的容量，地址就会改变
	bsh =(*reflect.SliceHeader)(unsafe.Pointer(&bss))
	p =unsafe.Pointer(bsh.Data)
	fmt.Println(unsafe.Pointer(&bss),bss,*bsh,p)

	s :=bytesToString(bs)
	sh = (*reflect.StringHeader)(unsafe.Pointer(&s))
	p = unsafe.Pointer(sh.Data)
	fmt.Println(&s,s,sh, p)
	fmt.Println(str)
}

func Test_Copy(t *testing.T)  {
	str:=strings.Repeat("abc",3)
	sh := (*reflect.StringHeader)(unsafe.Pointer(&str))
	p := unsafe.Pointer(sh.Data)
	fmt.Println(&str, str,  sh, p)



	bs:=[]byte(str)
	bsh :=(*reflect.SliceHeader)(unsafe.Pointer(&bs))
	p =unsafe.Pointer(bsh.Data)
	fmt.Println(unsafe.Pointer(&bs),bs,*bsh,p)

	bss :=bs[:]
	//copy(bss,[]byte(strings.Repeat("h",17)))//覆盖原底层数组内容，底层数组不允许被修改，地址不会改变，长度也不会改变
	bss=append(bss,strings.Repeat("h",8)...)	//添加的数据超过底层数组的容量，地址就会改变
	bsh =(*reflect.SliceHeader)(unsafe.Pointer(&bss))
	p =unsafe.Pointer(bsh.Data)
	fmt.Println(unsafe.Pointer(&bss),string(bss),string(bs),*bsh,p)



	s :=string(bs)
	sh = (*reflect.StringHeader)(unsafe.Pointer(&s))
	p = unsafe.Pointer(sh.Data)
	fmt.Println(&s,s,sh, p)
}