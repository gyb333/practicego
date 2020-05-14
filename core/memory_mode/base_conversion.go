package memory_mode

import "fmt"

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
	s := fmt.Sprintf("%b,%p", i,&i)			//格式化字符串
	fmt.Println(s)

}
