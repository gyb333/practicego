package main

import (
	"context"
	"go-grpc-example/pkg/gtls"
	pb "go-grpc-example/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

func main() {
	Address := "127.0.0.1:50052"
	tlsClient := gtls.Client{
		ServerName: "gyb",
		CertFile:   "conf/server.pem",
	}
	c, err := tlsClient.GetTLSCredentials()
	if err != nil {
		log.Fatalf("tlsClient.GetTLSCredentials err: %v", err)
	}

	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(c))
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()


	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(6 * time.Second)))
	defer cancel()


	client := pb.NewSearchServiceClient(conn)
	resp, err := client.Search(ctx, &pb.SearchRequest{
		Request: "gRPC",
	})
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Fatalln("client.Search err: deadline")
			}
		}
		log.Fatalf("client.Search err: %v", err)
	}

	log.Printf("resp: %s", resp.GetResponse())
}