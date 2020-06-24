package redis_util

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
)

type RedisLock struct {
	resource string
	token    string
	redis.Conn
	timeout  int
}

func (lock *RedisLock) tryLock() (ok bool, err error) {
	_, err = redis.String(
		lock.Do("SET", lock.key(), lock.token, "EX", int(lock.timeout), "NX"))
	if err == redis.ErrNil {
		// The lock was not successful, it already exists.
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (lock *RedisLock) Unlock() (err error) {
	_, err = lock.Do("del", lock.key())
	return
}

func (lock *RedisLock) key() string {
	return fmt.Sprintf("redislock:%s", lock.resource)
}

func (lock *RedisLock) AddTimeout(ex_time int64) (ok bool, err error) {
	ttl_time, err := redis.Int64(lock.Do("TTL", lock.key()))
	fmt.Println(ttl_time)
	if err != nil {
		log.Fatal("redis get failed:", err)
	}
	if ttl_time > 0 {
		fmt.Println(11)
		_, err := redis.String(lock.Do("SET", lock.key(), lock.token, "EX", int(ttl_time+ex_time)))
		if err == redis.ErrNil {
			return false, nil
		}
		if err != nil {
			return false, err
		}
	}
	return false, nil
}








