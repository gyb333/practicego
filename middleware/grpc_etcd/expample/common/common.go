package common

import "flag"

var (
	Serv = flag.String("service", "hello_service", "service name")
	Port = flag.Int("port", 50001, "listening port")
	Reg  = flag.String("reg", "http://172.23.0.20:2379", "register etcd address")
)
