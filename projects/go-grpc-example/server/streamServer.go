package main

import (
	"go-grpc-example/logging"
	pb "go-grpc-example/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"io"
	"log"
	"net"
)

type StreamService struct{}

func main() {
	Address := "127.0.0.1:50052"

	var logger = logging.NewLogger("grpcLogger")
	grpclog.SetLoggerV2(logging.NewZapLogger(logger))
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	// 实例化grpc Server
	s := grpc.NewServer()

	// 注册HelloService
	//pb.RegisterHelloServer(s, HelloService)
	pb.RegisterStreamServer(s,&StreamService{})
	grpclog.Infoln("Listen on " + Address)

	s.Serve(listen)
}
/*
服务器端流式 RPC
 */
func (s *StreamService) List(r *pb.StreamRequest, stream pb.Stream_ListServer) error {
	for n := 0; n <= 6; n++ {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  r.Pt.Name,
				Value: r.Pt.Value + int32(n),
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

/*
客户端流式 RPC
 */
func (s *StreamService) Record(stream pb.Stream_RecordServer) error {
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.StreamResponse{Pt: &pb.StreamPoint{Name: "gRPC Stream Server: Record", Value: 1}})
		}
		if err != nil {
			return err
		}

		log.Printf("stream.Recv pt.name: %s, pt.value: %d", r.Pt.Name, r.Pt.Value)
	}

	return nil
}

/*
双向流式 RPC
 */
func (s *StreamService) Route(stream pb.Stream_RouteServer) error {
	n := 0
	for {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  "gPRC Stream Client: Route",
				Value: int32(n),
			},
		})
		if err != nil {
			return err
		}

		r, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		n++

		log.Printf("stream.Recv pt.name: %s, pt.value: %d", r.Pt.Name, r.Pt.Value)
	}

	return nil
}