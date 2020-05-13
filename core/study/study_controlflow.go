package basic

import "fmt"

//全局变量的声明
var (
	i      = 1
	strSQL = `
	多行字符串声明
`
)


func ControlFlow()  {
	for j := 0; j < 10; j++ {
		fmt.Println(j)
	}
	for {
		fmt.Println("相当于While(True)")
		break
	}
	for i < 5 {
		fmt.Println("相当于While(条件)", i)
		i++
	}

	sw := 1
	switch sw {
	case 0:
		fmt.Println("sw=0")
	case 1:
		fmt.Println("sw=1")
	}
	switch {
	case sw == 0:
		fmt.Println("sw条件=0")
		fallthrough //继续执行下一个case 不需要判断case 条件
	case sw == 1:
		fmt.Println("sw条件=1")
	}
	switch sw := 0; {
	case sw == 0:
		fmt.Println("sw初始化=0", sw)
		fallthrough //继续执行下一个case 不需要判断case 条件
	case sw == 1:
		fmt.Println("sw初始化=1", sw)
	case sw == 2:
		fmt.Println("sw初始化=2", sw)
	}

LABEL:
	for {
		for j := 0; j < 10; j++ {
			if j > 2 {
				break LABEL
			} else {
				fmt.Println("break LABEL 跳出多层循环", j)
			}
		}
	}

CLABEL:
	for j := 0; j < 5; j++ {
		for {
			fmt.Println("continue CLABEL 继续标签层循环", j)
			continue CLABEL
		}
	}
}
