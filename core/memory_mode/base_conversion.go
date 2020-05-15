package memory_mode

import (
	"encoding/binary"
	"fmt"
	"math"
	"unsafe"
)

/*

%b     二进制表示
%c     相应Unicode码点所表示的字符
%d     十进制表示
%o     八进制表示
%q     单引号围绕的字符字面值，由Go语法安全地转义
%x     十六进制表示，字母形式为小写 a-f
%X     十六进制表示，字母形式为大写 A-F
%U     Unicode格式：U+1234，等同于 "U+%04X"
%p(16位表示法,前导0x)
 */
func BaseConversion()  {
	i:=55
	fmt.Printf("%T,%v,%b,%o,%X\n",i,i,i,i,i)
	pb :=  unsafe.Pointer(uintptr(i))		//强制转为指针类型
	fmt.Printf("%x,%p,%T\n",i,pb,pb)
	s := fmt.Sprintf("%b,%p", i,&i)			//格式化字符串
	fmt.Println(s)

	var f float32
	var ip int32

	// unsafe
	f = 1.234
	ip = *((*int32)(unsafe.Pointer(&f)))
	fmt.Printf("%f,%x,%d,%x\n",f,f, ip,ip)

	// safe
	var tmp [4]byte
	f = 1.234
	fmt.Printf("%f,%x\n",f,f)
	ip=int32(f)
	fmt.Printf("%d,%x\n",ip,ip)
	binary.LittleEndian.PutUint32(tmp[:], math.Float32bits(f))
	ip = int32(binary.LittleEndian.Uint32(tmp[:]))
	fmt.Printf("%d,%x\n",ip,ip)

}
