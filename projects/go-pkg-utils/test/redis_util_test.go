package test

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"go-pkg-utils/redis_util"
	"go-pkg-utils/setting_util"
	"log"
	"os"
	"testing"
)

func init() {
	fmt.Println(os.Getwd())
	setting := setting_util.NewSetting("../conf/config.ini")
	redis_util.Setup(setting)
}

func TestRedisPipelining(t *testing.T) {
	c := redis_util.RedisPool.Get()
	c.Send("SET", "key", "bar")
	c.Send("GET", "key")
	c.Flush()
	c.Send("SET", "kv", "kv")
	c.Send("GET", "kv")
	c.Receive()
	v, err := c.Receive()
	if err != nil {
		log.Println(err)
	}
	log.Println(v)

	c.Send("MULTI")
	c.Send("INCR", "foo")
	c.Send("INCR", "bar")
	v, err = c.Do("EXEC")
	if err != nil {
		log.Println(err)
	}
	log.Println(v)
}

//构造实际场景的hash结构体
var p1, p2 struct {
	Description string `redis:"description"`
	Url         string `redis:"url"`
	Author      string `redis:"author"`
}

func errCheck(tp string,err error) {
	if err != nil {
		fmt.Printf("sorry,has some error for %s %v.\r\n",tp,err)
		os.Exit(-1)
	}
}

func TestRedisDoArgs(t *testing.T) {
	c := redis_util.RedisPool.Get()
	p1.Description = "my blog"
	p1.Url = "http://xxbandy.github.io"
	p1.Author = "bgbiao"
	reply, err := c.Do("hmset", redis.Args{}.Add("hao123").AddFlat(&p1)...)
	errCheck("hmset",err)
	log.Println(reply)
	m := map[string]string{
		"description": "oschina",
		"url":         "http://my.oschina.net/myblog",
		"author":      "xxbandy",
	}

	reply, err = c.Do("hmset", redis.Args{}.Add("hao").AddFlat(m)...)
	errCheck("hmset",err)
	log.Println(reply)

	for _, key := range []string{"hao123", "hao"} {
		//等同于hgetall的输出类型，输出字符串为k/v类型
		hashV, _ := redis.StringMap(c.Do("hgetall", key))
		fmt.Println("hgetall ",hashV)
		//等同于hmget 的输出类型，输出字符串到一个字符串列表
		hashV2, _ := redis.Strings(c.Do("hmget", key, "description", "url", "author"))
		for _, hashv := range hashV2 {
			fmt.Println("hmget: ",hashv)
		}

		v, _ := redis.Values(c.Do("hgetall", key))
		if err := redis.ScanStruct(v, &p2); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("redis.ScanStruct %+v\n", p2)
	}


}
func TestRedisDoString(t *testing.T) {
	c := redis_util.RedisPool.Get()
	c.Do("set","name","set")
	c.Do("set","id",11)
	name,err := redis.String(c.Do("get","name"))
	errCheck("get",err)
	id,err := redis.Int(c.Do("get","id"))
	errCheck("get",err)
	fmt.Println("book name:",name,id)
}
func TestRedisDoHash(t *testing.T) {
	c := redis_util.RedisPool.Get()
	_,err := c.Do("hset","books","name","golang")
	errCheck("hset",err)

	r,err := redis.String(c.Do("hget","books","name"))
	errCheck("hset",err)
	fmt.Println("book name:",r)

}



func TestRedisDoSet(t *testing.T) {
	c := redis_util.RedisPool.Get()
	c.Do("mset","id",100,"fn",200)
	c.Do("expire","id",100)

	if r,err := redis.Ints(c.Do("mget","id","fn")); err == nil {
		for _,v := range r{
			fmt.Println("value:",v)
		}
	}
}

func TestRedisDoList(t *testing.T) {
	c := redis_util.RedisPool.Get()
	c.Do("lpush","updateid","0407","0408","0409")

	r,err := redis.String(c.Do("lpop","updateid"))
	errCheck("lpop",err)
	log.Println(r)
}