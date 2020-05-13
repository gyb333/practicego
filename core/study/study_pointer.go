package basic

import "fmt"

func StudyPointer()  {
	pointerBase()
}

func pointerBase()  {
	var pi *int
	fmt.Printf("表示指向int类型的指针%T,%v\n",pi,pi)
	pi = new(int)
	fmt.Println(pi,*pi)

	var ap [5]*int
	fmt.Printf("表示长度为5,int类型的指针数组%T,%v\n",ap,ap)

	var pa  *[5]int
	fmt.Printf("表示指向长度为5,int类型的数组指针%T,%v\n",pa,pa)

	var ppa **[5]int
	fmt.Printf("表示指向长度为5,int类型的数组指针的指针%T,%v\n",ppa,ppa)

	var pap *[5]*int
	fmt.Printf("表示指向长度为5,int类型指针的数组指针%T,%v\n",pap,pap)


}

