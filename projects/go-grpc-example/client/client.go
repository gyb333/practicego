package main

import (
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	zipkinot  "github.com/openzipkin-contrib/zipkin-go-opentracing"
	zipkin "github.com/openzipkin/zipkin-go"

	"go-grpc-example/logging"
	"go-grpc-example/pkg/gtls"
	pb "go-grpc-example/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
	SERVICE_NAME              = "zipkin_client"
	ZIPKIN_HTTP_ENDPOINT      = "http://hadoop:9411/api/v1/spans"
	ZIPKIN_RECORDER_HOST_PORT = "127.0.0.1:9000"
)

func main() {
	var logger = logging.NewLogger("grpcLogger")
	grpclog.SetLoggerV2(logging.NewZapLogger(logger))

	reporter :=  zipkinhttp.NewReporter(ZIPKIN_HTTP_ENDPOINT)
	defer reporter.Close()
	endpoint, err := zipkin.NewEndpoint(SERVICE_NAME, ZIPKIN_RECORDER_HOST_PORT)
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}
	// initialize our tracer
	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}
	// use zipkin-go-opentracing to wrap our tracer
	tracer := zipkinot.Wrap(nativeTracer)

	// optionally set as Global OpenTracing tracer instance
	opentracing.SetGlobalTracer(tracer)




	//tlsClient := gtls.Client{
	//	ServerName: "gyb",
	//	CertFile:   "conf/server.pem",
	//}
	//tc, err := tlsClient.GetTLSCredentials()
	//if err != nil {
	//	log.Fatalf("tlsClient.GetTLSCredentials err: %v", err)
	//}
	tlsClient := gtls.Client{
		ServerName: "gyb333",
		CaFile:     "conf/CA/ca.pem",
		CertFile:   "conf/CA/clientCA.pem",
		KeyFile:    "conf/CA/client.key",
	}

	tc, err := tlsClient.GetCredentialsByCA()
	if err != nil {
		log.Fatalf("GetTLSCredentialsByCA err: %v", err)
	}

	auth := Auth{
		AppKey:    "gyb333",
		AppSecret: "20200101",
	}

	opts :=[]grpc.DialOption{
		grpc.WithTransportCredentials(tc),
		grpc.WithPerRPCCredentials(&auth),
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(tracer, otgrpc.LogPayloads()),
		),
	}
	// 连接
	//conn, err := grpc.Dial(Address, grpc.WithInsecure())
	conn, err := grpc.Dial(Address, opts...)

	if err != nil {
		grpclog.Fatalln(err)
	}

	defer conn.Close()

	// 初始化客户端
	c := pb.NewHelloClient(conn)

	// 调用方法
	reqBody := new(pb.HelloRequest)
	reqBody.Name = "gRPC ZipKin"
	r, err := c.SayHello(context.Background(), reqBody)
	if err != nil {
		grpclog.Fatalln(err)
	}

	grpclog.Infoln(r.Message)
}



type Auth struct {
	AppKey    string
	AppSecret string
}

func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_key": a.AppKey, "app_secret": a.AppSecret}, nil
}

func (a *Auth) RequireTransportSecurity() bool {
	return true
}