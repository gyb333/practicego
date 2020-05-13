package basic

import (
	"fmt"
	"unsafe"
			"reflect"
	"bytes"
)

func StudySlice()  {
	sliceBase()
	sliceDataType()
}

func sliceBase()  {
	a := []byte("a,b,c")
	b := bytes.Split(a, []byte(","))
	fmt.Println(a,b)
	a =[]byte(str)
	b = bytes.Split(a,[]byte(""))
	fmt.Println(string(a),b)
}

func sliceDataType()  {
	var bs []byte
	fmt.Printf("%p,%T,%d,%v,%d,%d\n",&bs,bs,unsafe.Sizeof(bs),bs,len(bs),cap(bs))
	sh :=(*reflect.SliceHeader)(unsafe.Pointer(&bs))
	p :=unsafe.Pointer(sh.Data)
	fmt.Println(*sh,sh.Data,sh.Len,sh.Cap,p,bs)

	bs =[]byte{}
	p =unsafe.Pointer(sh.Data)
	fmt.Println(*sh,sh.Data,sh.Len,sh.Cap,p,bs,"-"+string(bs)+"-")

	bs = []byte(str)
	p =unsafe.Pointer(sh.Data)
	fmt.Println(*sh,sh.Data,sh.Len,sh.Cap,p,bs,string(bs))

	bArray :=(*[15]byte)(p)
	fmt.Printf("%p,%T,%v,%s,%c\n", unsafe.Pointer(bArray),bArray, *bArray, *bArray, bArray[0])
	for i,v:=range bArray{
		if i==0{
			bArray[i]=v+32		//切片底層數組內容可以修改
		}
	}
	fmt.Println(bs,string(bs))
	fmt.Println(bArray,string(bArray[:]))

	cs :=make([]byte,sh.Len)
	for i,v :=range bs{
		if i==0{
			cs[i] =v-32
		}else {
			cs[i] =v
		}
	}
	fmt.Println(cs,string(cs))
	cs =make([]byte,sh.Len)
	for i,_ :=range bs{
		cs[i]= *(*byte)(unsafe.Pointer(uintptr(p) + uintptr(i)))
	}
	fmt.Println(cs,string(cs))


}







