package basic

import (
	"fmt"
	"unsafe"
	"reflect"
)

func NewMake()  {
	fmt.Println("------------------var arr []int直接引用------------------")
	var arr []int
	fmt.Println(unsafe.Pointer(&arr),reflect.TypeOf(arr).String(),arr,len(arr),cap(arr))
	ap := (*reflect.SliceHeader)(unsafe.Pointer(&arr))
	fmt.Println(ap,ap.Data,ap.Len,ap.Cap)
	p := (*int)(unsafe.Pointer(ap.Data))
	fmt.Println(p,p==nil)
	fmt.Println("------------------arr =[]int{}直接引用------------------")
	arr =[]int{}
	fmt.Println(unsafe.Pointer(&arr),reflect.TypeOf(arr).String(),arr,len(arr),cap(arr))
	ap = (*reflect.SliceHeader)(unsafe.Pointer(&arr))
	fmt.Println(ap,ap.Data,ap.Len,ap.Cap)
	p = (*int)(unsafe.Pointer(ap.Data))
	fmt.Println(p,p==nil)

	fmt.Println("------------------arr =make([]int,0)直接引用------------------")
	arr =make([]int,0)
	fmt.Println(unsafe.Pointer(&arr),reflect.TypeOf(arr).String(),arr,len(arr),cap(arr))
	ap = (*reflect.SliceHeader)(unsafe.Pointer(&arr))
	fmt.Println(ap,ap.Data,ap.Len,ap.Cap)
	p = (*int)(unsafe.Pointer(ap.Data))
	fmt.Println(p,p==nil)

	fmt.Println("------------------arrp := new ([]int)间接引用------------------")

	arrp := new ([]int)
	fmt.Println(unsafe.Pointer(arrp),reflect.TypeOf(arrp).String(),arrp,len(*arrp),cap(*arrp))
	ap = (*reflect.SliceHeader)(unsafe.Pointer(arrp))
	fmt.Println(ap,ap.Data,ap.Len,ap.Cap)
	p = (*int)(unsafe.Pointer(ap.Data))
	fmt.Println(p,p==nil)

	fmt.Println("------------------arrp = &[]int{}间接引用------------------")
	arrp = &[]int{}
	fmt.Println(unsafe.Pointer(arrp),reflect.TypeOf(arrp).String(),arrp,len(*arrp),cap(*arrp))
	ap = (*reflect.SliceHeader)(unsafe.Pointer(arrp))
	fmt.Println(ap,ap.Data,ap.Len,ap.Cap)
	p = (*int)(unsafe.Pointer(ap.Data))
	fmt.Println(p,p==nil)


}
