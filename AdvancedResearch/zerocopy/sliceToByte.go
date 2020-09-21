package zerocopy

import (
	"reflect"
	"unsafe"
)

func stringToBytes(s string) []byte  {
	x :=(*[2]uintptr)(unsafe.Pointer(&s))
	h :=[3]uintptr{x[0],x[1],x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
func stringToBytes2(s string)[]byte  {
	ptr :=(*[2]uintptr)(unsafe.Pointer(&s))
	sh := struct {
		addr uintptr
		len int
		cap int
	}{ptr[0], len(s), len(s)}

	return  *(*[]byte)(unsafe.Pointer(&sh));
}

func stringToByteSlice(s string)[]byte  {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: stringHeader.Data,
		Len:  stringHeader.Len,
		Cap:  stringHeader.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func stringToByteSlice2(s string)[]byte  {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	var bs []byte
	sh :=(*reflect.SliceHeader)(unsafe.Pointer(&bs))
	sh.Cap=stringHeader.Len
	sh.Len=stringHeader.Len
	sh.Data=stringHeader.Data
	return bs;
}


func bytesToString(b []byte) string  {
	return *(*string)(unsafe.Pointer(&b))
}