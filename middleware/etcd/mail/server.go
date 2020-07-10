package main

import (
	"GYB.Middleware/etcd"
	"GYB.Middleware/etcd/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"

)

type service struct {
}

func (s *service) SendMail(ctx context.Context, req *proto.MailRequest) (res *proto.MailResponse, err error) {
	fmt.Printf("邮箱:%s;发送内容:%s", req.Mail, req.Text)
	return &proto.MailResponse{
		Ok: true,
	}, nil
}

func main() {

	// 监听本地的8972端口
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()                       // 创建gRPC服务器

	proto.RegisterMailServiceServer(s, &service{}) // 在gRPC服务端注册服务

	reflection.Register(s) //在给定的gRPC服务器上注册服务器反射服务

	// Serve方法在lis上接受传入连接，为每个连接创建一个ServerTransport和server的goroutine。
	// 该goroutine读取gRPC请求，然后调用已注册的处理程序来响应它们。

	//etcd服务注册
	reg, err := etcd.NewService(etcd.ServiceInfo{
		Name: "g.srv.mail",
		IP:   "127.0.0.1:8972", //grpc服务节点ip
	}, []string{"127.23.0.20:2379", "127.23.0.21:2379", "127.23.0.22:2379",}) // etcd的节点ip
	if err != nil {
		log.Fatal(err)
	}
	go reg.Start()

	if err := s.Serve(lis); err != nil {
		fmt.Println(err)
	}
}
