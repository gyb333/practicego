package Redis

import (
	"github.com/garyburd/redigo/redis"
	. "middleware/RedisRedis/common"
)

func AddZset(payloadKey string, score int64, zsetName string) (err error) {
	con := Pool.Get()
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
	con := Pool.Get()
	defer con.Close()
	_, err = con.Do("zrem", redis.Args{}.Add(zsetName).AddFlat(payloadKeys)...) //TODO 这个点易错。
	return
}

//index sorted set from start to end, [start:end], eg: [0:1] will return[member1, score1, member2, score2]
func RangeZset(start, end int, zsetName string) (payloadKeys []string, err error) {
	con := Pool.Get()
	defer con.Close()
	payloadKeys, err = redis.Strings(con.Do("zrange", zsetName, start, end, "withscores"))
	return
}

func RangeZsetByScore(start, end int64, zsetName string) (payloadKeys []string, err error) {
	con := Pool.Get()
	defer con.Close()
	payloadKeys, err = redis.Strings(con.Do("zrangebyscore",zsetName, start, end, "withscores"))
	return
}
