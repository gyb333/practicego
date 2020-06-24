package redis_util

import (
	"go-pkg-utils/setting_util"
	"time"
)

type RedisSetting struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}
var rs =&RedisSetting{}

func initRedis(cfg *setting_util.Setting){
	 cfg.MapTo("redis",rs)
	 rs.IdleTimeout = rs.IdleTimeout * time.Second
}
