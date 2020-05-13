package main

import (
	"./sort"
	"fmt"
)

func main() {
	var arr=[]int{3,5,6,8,9,4,1,2,7}
	sort.BubbleSort(arr)
	fmt.Println(arr)
}
