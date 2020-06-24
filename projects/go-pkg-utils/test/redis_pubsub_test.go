package test

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"go-pkg-utils/redis_util"
	"go-pkg-utils/setting_util"
	"testing"
	"time"
)

//订阅
func subs(forever chan<- bool)  {
	//订阅
	conn := redis_util.RedisPool.Get()
	defer conn.Close()

	psc := redis.PubSubConn{conn}
	_ = psc.Subscribe("chan_go")
	for{
		switch v := psc.Receive().(type) {
		case redis.Message:
			fmt.Printf("接收消息 %s: message: %s\n", v.Channel, v.Data)
			forever<-true
		case redis.Subscription:
			fmt.Printf("发布消息 %s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			fmt.Println(v)
			return
		}
	}
}

//发布消息
func push(message string)  {
	//发布消息
	conn := redis_util.RedisPool.Get()
	defer conn.Close()


	_,err := conn.Do("Publish","chan_go",message)
	if err != nil{
		fmt.Println(err)
		return
	}
}

func init() {
	setting := setting_util.NewSetting("../conf/config.ini")
	redis_util.Setup(setting)
}

func TestRedisPubSub(t *testing.T) {

	forever:=make(chan bool)
	go subs(forever)
	time.Sleep(time.Second)
	push("this is wd")
	<-forever
}
