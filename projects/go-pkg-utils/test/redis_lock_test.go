package test

import (
	"fmt"
	"go-pkg-utils/redis_util"
	"go-pkg-utils/setting_util"
	"log"
	"testing"
	"time"
)

func init() {
	setting := setting_util.NewSetting("../conf/config.ini")
	redis_util.Setup(setting)
}
/*
$ go run lock.go
start
10
11
end

如果同时起多个进程去测试，会遇到这么一个结果:
$ go run lock.go
start
2016/03/23 01:23:22 Lock
exit status 1
*/
func TestRedisLock(t *testing.T) {
	fmt.Println("start")
	DefaultTimeout := 10
	lock, ok, err := redis_util.TryLock("xiaoru.cc", "token", int(DefaultTimeout))
	if err != nil {
		log.Fatal("Error while attempting lock")
	}
	if !ok {
		log.Fatal("Lock")
	}
	lock.AddTimeout(100)

	time.Sleep(time.Duration(DefaultTimeout) * time.Second)
	fmt.Println("end")
	defer lock.Unlock()
}
