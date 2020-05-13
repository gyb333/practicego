package basic

import (
	"fmt"
	)

/*
继承：通过匿名结构体的嵌套来实现is a
      使用组合的方式实现 has a
覆盖：子类可以重新实现父类的方法

多态：把接口作为参数或者返回值类型传递
	接口可以用任何实现该接口的指针来实例化


方法：某个类型的行为功能，需要知道接受者调用
	方法名可以相同，只要接受者不同就行
函数：一段独立的代码,可以直接调用
	函数在同一个包下不能冲突
 */




type Person struct {
	name string
	age int
}
type Student struct {
	Person		//匿名字段默认为类型名称
	s School
}

type School struct{
	name string
	address string
}

type ManyState interface {
	PersonName() string
	SchoolName() string
}

//匿名字段对应方法继承
func (p *Person) PersonName() string {
	return "Person"+p.name
}

func (s *School) SchoolName() string {
	return "School"+s.name
}

//匿名字段对应方法重写
func (s *Student) PersonName() string {
	return "Student"+s.name
}

func (s *Student) SchoolName() string {
	return "Student"+s.s.name
}
//函数
func OOPMain()  {
	p:=Person{name:"张三",age:10}
	fmt.Println(p.name,p.age)
	s:=Student{Person:p,s:School{"学校","地址"}}
	fmt.Println(s.name,s.age,s.s.name,s.s.address)
	fmt.Println(s.PersonName(),s.Person.PersonName())

	fmt.Println(s.SchoolName(),s.s.SchoolName())

	var ms ManyState	//接口类型
	ms =&s
	fmt.Println("接口类型调用实现多态：",ms.PersonName(),ms.SchoolName())

	//类型断言
	if i,ok :=ms.(*Student); ok{
		fmt.Printf("%p,%p,%v\n",ms,i,i)
	}
}

