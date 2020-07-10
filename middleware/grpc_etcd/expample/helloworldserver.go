package main

import (
	grpclb "GYB.Middleware/grpc_etcd"
	"GYB.Middleware/grpc_etcd/expample/common"
	"GYB.Middleware/grpc_etcd/expample/pb"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"


	"golang.org/x/net/context"
	"google.golang.org/grpc"
)



func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *common.Port))
	if err != nil {
		panic(err)
	}

	err = grpclb.Register(*common.Serv, "127.0.0.1", *common.Port, *common.Reg, time.Second*10, 15)
	if err != nil {
		panic(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		log.Printf("receive signal '%v'", s)
		grpclb.UnRegister()
		os.Exit(1)
	}()

	log.Printf("starting hello service at %d", *common.Port)
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	s.Serve(lis)
}

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Printf("%v: Receive is %s\n", time.Now(), in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}
