package main

import (
	"github.com/garyburd/redigo/redis"
	"fmt"
	. "middleware/RedisRedis/common"
	"time"
)

//订阅
func subs(forever chan<- bool)  {
	//订阅
	conn := Pool.Get()
	defer conn.Close()

	fmt.Println("接收消息....")

	psc := redis.PubSubConn{conn}
	_ = psc.Subscribe("chan_go")
	for{
		switch v := psc.Receive().(type) {
		case redis.Message:
			fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
			forever<-true
		case redis.Subscription:
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			fmt.Println(v)
			return
		}
	}
}

//发布消息
func push(message string)  {
	//发布消息
	conn := Pool.Get()
	defer conn.Close()

	fmt.Println("发布消息....")

	_,err := conn.Do("Publish","chan_go",message)
	if err != nil{
		fmt.Println(err)
		return
	}
}


func main() {
	Init("tcp","hadoop:6379")
	forever:=make(chan bool)
	go subs(forever)
	time.Sleep(time.Second)
	push("this is wd")
	<-forever
}
