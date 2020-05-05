package common

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"fmt"
)

type RedisKv struct {
	Key   string
	Value string
	TTR   int64 //time to return
}

var Pool *redis.Pool

//建立连接池
func Init(network, address string)  {
	Pool = &redis.Pool{
		MaxIdle:     8,
		IdleTimeout: 100,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(network, address)
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}
	conn := Pool.Get()
	defer conn.Close()
	if _, err := conn.Do("ping"); err != nil {
		panic(err)
	}
}


func FailOnError(err error, msg string){
	if err != nil {
		log.Fatalf("%s:%s", msg, err)
		panic(fmt.Sprintf("%s:%s", msg, err))
	}
}

func GetRedisConn(network, address string)  redis.Conn{

	conn,err:=redis.Dial(network,address)
	FailOnError(err,"Failed to connect to Redis!")
	return conn
}

//实现简单队列
func BatchPushQueue(queueName string, keys []string) (err error) {
	if len(keys) == 0 {
		return
	}
	con := Pool.Get()
	defer con.Close()
	_, err = con.Do("lpush", redis.Args{}.Add(queueName).AddFlat(keys)...)
	return
}

func PopQueue(queueName string, timeout int) (data string, err error) {
	con := Pool.Get()
	defer con.Close()
	nameAndData, err := redis.Strings(con.Do("brpop", queueName, timeout))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		return
	}
	if len(nameAndData) > 1 {
		data = nameAndData[1]
	}
	return
}

