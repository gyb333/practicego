package main

import (
	"GYB.Middleware/grpc_etcd/expample/common"
	"GYB.Middleware/grpc_etcd/expample/pb"
	"flag"
	"fmt"
	"time"

	"strconv"

	grpclb "GYB.Middleware/grpc_etcd"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//var (
//	serv = flag.String("service", "hello_service", "service name")
//	reg  = flag.String("reg", "http://127.0.0.1:2379", "register etcd address")
//)

func main() {
	flag.Parse()
	fmt.Println("serv", *common.Serv)
	r := grpclb.NewResolver(*common.Serv)
	b := grpc.RoundRobin(r)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	conn, err := grpc.DialContext(ctx, *common.Reg, grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		panic(err)
	}
	fmt.Println("conn...")

	ticker := time.NewTicker(1 * time.Second)
	for t := range ticker.C {
		client := pb.NewGreeterClient(conn)
		resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "world " + strconv.Itoa(t.Second())})
		if err == nil {
			fmt.Printf("%v: Reply is %s\n", t, resp.Message)
		} else {
			fmt.Println(err)
		}
	}
}
