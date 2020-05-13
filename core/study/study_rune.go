package basic

import (
	"fmt"
	"unsafe"
				"strconv"
	"unicode/utf8"
)
var r rune ='h'
var u ='世'

func StudyRune()  {
	runeType()
	runeSlice()
	byteRune()
}


func runeType()  {
	fmt.Printf("%p,%T,%d,%U,%v,%c\n",&r,r,unsafe.Sizeof(r),r,r,r)
	fmt.Printf("%p,%T,%d,%U,%v,%c\n",&u,u,unsafe.Sizeof(u),u,u,u)

	fmt.Println("--------------------------------------------")
	fmt.Printf("%d,%q,%#X\n",r,r,r)
	fmt.Printf("%d,%q,%#X\n",u,u,u)

	fmt.Println("--------------------------------------------")
	fmt.Printf("%c,%c,%c\n",r,104,0x68)
	fmt.Printf("%c,%c,%c\n",u,19990,0x4E16)

	fmt.Println("--------------------------------------------")
	i,_:=strconv.ParseInt("68",16,32)
	fmt.Printf("%d,%q,%X,%t,%t,%t\n",i,i,i,r==rune(i),104==i,104==0x68)
	i,_=strconv.ParseInt("4E16",16,32)
	fmt.Printf("%d,%q,%X,%t,%t,%t\n",i,i,i,u==rune(i),19990==i,19990==0x4E16)
}


func runeSlice()  {
	var s =[]byte(string('h'))
	fmt.Printf("%d,%d,%v\t%# X\n",len(s),cap(s),s,s)
	var us =[]byte(string('世'))
	fmt.Printf("%d,%d,%v\t%# X\n",len(us),cap(us),us,us)

	fmt.Println("--------------------------------------------")

	fmt.Printf("%s\t%s\t%s\t%s\n","h","\x68","\u0068","\U00000068")
	fmt.Printf("%s\t%s\t%s\t%s\n","世","\xE4\xB8\x96","\u4E16","\U00004E16")


}

func byteRune()  {
	s:=[]byte{228,184,150}
	fmt.Printf("%d,%d,%v\t%# X\n",len(s),cap(s),s,s)
	dr,_:=utf8.DecodeRune(s)
	fmt.Printf("%c\t%t\n",dr,dr==u)
	s=[]byte{0xE4,0xB8,0x96}
	fmt.Printf("%d,%d,%v\t%# X\n",len(s),cap(s),s,s)
	dr,_=utf8.DecodeRune(s)
	fmt.Printf("%c\t%t\n",dr,dr==u)
}