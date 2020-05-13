package basic

import (
	"regexp"
	"fmt"
)

const text ="My Email is 8768840@qq.com or zhongduzhi@163.com or gyb333@126.com.cn"

func StudyRegex()  {
	re :=regexp.MustCompile(`([a-zA-Z0-9]+)@([a-zA-Z0-9]+)\.([a-zA-Z0-9.]+)`)
	match :=re.FindAllString(text,-1)
	 fmt.Println(match)
	smatch :=re.FindAllStringSubmatch(text,-1)
	fmt.Println(smatch)
}
