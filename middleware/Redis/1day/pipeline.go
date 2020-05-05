package main

import (
	. "../common"
	"fmt"
	"strconv"
	"time"
)

func main() {
	c:= GetRedisConn("tcp","docker:6379")
	defer c.Close()
	start:=time.Now()
	for j:=0;j<1000;j++{
		c.Do("HSET","hashkey:"+strconv.Itoa(j),
			 	"field"+strconv.Itoa(j),"value"+strconv.Itoa(j))
	}
	end := time.Now()
	fmt.Println(end.Sub(start))
	for j:=0;j<10000;j++{
		c.Send("HSET","mhashkey:"+strconv.Itoa(j),
				"field"+strconv.Itoa(j),"value"+strconv.Itoa(j))
		c.Flush()
	}
	t := time.Now()
	fmt.Println(t.Sub(end))

	for j:=0;j<1000000;j++{
		c.Send("pfadd","Hyperloglog","key"+strconv.Itoa(j))
		c.Flush()
	}
	t1 := time.Now()
	fmt.Println(t1.Sub(t))
}
