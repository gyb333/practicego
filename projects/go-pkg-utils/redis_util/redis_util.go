package redis_util

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"go-pkg-utils/setting_util"
	"time"
)

var RedisPool *redis.Pool

func Setup(cfg *setting_util.Setting) error {
	initRedis(cfg)
	RedisPool = &redis.Pool{
		MaxIdle:     rs.MaxIdle,
		MaxActive:   rs.MaxActive,
		IdleTimeout: rs.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c,err := redis.Dial("tcp",rs.Host,
				redis.DialKeepAlive(1*time.Second),
				redis.DialPassword(rs.Password),
				redis.DialConnectTimeout(5*time.Second),
				redis.DialReadTimeout(1*time.Second),
				redis.DialWriteTimeout(1*time.Second))
			//c, err := redis.Dial("tcp", rs.Host)
			if err != nil {
				return nil, err
			}
			//if rs.Password != "" {
			//	if _, err := c.Do("AUTH", rs.Password); err != nil {
			//		c.Close()
			//		return nil, err
			//	}
			//}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return nil
}

func Set(key string, data interface{}, time int) (bool, error) {
	conn := RedisPool.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	reply, err := redis.Bool(conn.Do("SET", key, value))
	conn.Do("EXPIRE", key, time)

	return reply, err
}

func Exists(key string) bool {
	conn := RedisPool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

func Get(key string) ([]byte, error) {
	conn := RedisPool.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func Delete(key string) (bool, error) {
	conn := RedisPool.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

func LikeDeletes(key string) error {
	conn := RedisPool.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}

//实现简单队列
func PushQueue(queueName string, key string) (err error) {
	con := RedisPool.Get()
	defer con.Close()
	_, err = con.Do("lpush", queueName, key)
	if err != nil {
		return
	}
	return
}

func BatchPushQueue(queueName string, keys []string) (err error) {
	if len(keys) == 0 {
		return
	}
	con := RedisPool.Get()
	defer con.Close()
	_, err = con.Do("lpush", redis.Args{}.Add(queueName).AddFlat(keys)...)
	return
}

func PopQueue(queueName string, timeout int) (data string, err error) {
	con := RedisPool.Get()
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

func TryLock( resource string, token string, DefaulTimeout int) (lock *RedisLock, ok bool, err error) {
	return TryLockWithTimeout(resource, token, DefaulTimeout)
}

func TryLockWithTimeout( resource string, token string, timeout int) (lock *RedisLock, ok bool, err error) {
	lock = &RedisLock{resource, token, RedisPool.Get(), timeout}

	ok, err = lock.tryLock()

	if !ok || err != nil {
		lock = nil
	}

	return
}


func AddZset(payloadKey string, score int64, zsetName string) (err error) {
	con := RedisPool.Get()
	defer con.Close()
	_, err = con.Do("zadd", zsetName, score, payloadKey)
	if err != nil {
		return
	}
	return
}

func RemZset(payloadKeys []string, zsetName string) (err error) {
	if len(payloadKeys) == 0 {
		return
	}
	con := RedisPool.Get()
	defer con.Close()
	_, err = con.Do("zrem", redis.Args{}.Add(zsetName).AddFlat(payloadKeys)...) //TODO 这个点易错。
	return
}

//index sorted set from start to end, [start:end], eg: [0:1] will return[member1, score1, member2, score2]
func RangeZset(start, end int, zsetName string) (payloadKeys []string, err error) {
	con := RedisPool.Get()
	defer con.Close()
	payloadKeys, err = redis.Strings(con.Do("zrange", zsetName, start, end, "withscores"))
	return
}

func RangeZsetByScore(start, end int64, zsetName string) (payloadKeys []string, err error) {
	con := RedisPool.Get()
	defer con.Close()
	payloadKeys, err = redis.Strings(con.Do("zrangebyscore",zsetName, start, end, "withscores"))
	return
}


