package basic

import (
	"fmt"
	)

func StudyMap()  {
	mapDataType()
}

func mapDataType()  {
	var m map[string]string
	fmt.Printf("%p,%T,%v,%t\n",&m,m,m,m==nil)

	m=make(map[string]string)
	m["key"]="value"

	m["gyb"]="gyb333"
	if _ ,ok :=m["gyb"];ok{
		delete(m,"gyb")
	}
	for k,v:=range m{
		fmt.Println(k,v)
	}





}