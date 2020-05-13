package basic

import "fmt"

/**
接口与类型之间是非嵌入式的
 */

func InterfaceMain()  {


	 baseInterface()
}


type Inter interface {
	Ping()
	Pang()
}

type St struct{}

//方法：限定了接受者
func (St) Ping() {
	fmt.Println("ping")
}
//方法：限定了接受者
func (*St) Pang() {
	fmt.Println("Pang")
}

func baseInterface()  {
	var st *St = nil
	var it Inter = st
	fmt.Printf("%T,%p,%#V,%d\n", st, st, st, st)
	fmt.Printf("%T,%p,%#V,%d\n", it, it, it, it)
	fmt.Println(st == it, st == nil, it == nil, st, it,&it)
	if it != nil {
		//it.Ping()
		it.Pang()
	}


	var i Inter = &St{}
	i.Pang()
}



