package main

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"

	"go-grpc-example/logging"
	"go-grpc-example/pkg/gtls"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"runtime/debug"

	pb "go-grpc-example/proto" // 引入编译生成的包

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	zipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/opentracing/opentracing-go"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)



const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
	SERVICE_NAME              = "zipkin_server"
	ZIPKIN_HTTP_ENDPOINT      = "http://hadoop:9411/api/v1/spans"
	ZIPKIN_RECORDER_HOST_PORT = "127.0.0.1:9000"
)

// 定义helloService并实现约定的接口
type helloService struct{
	auth *Auth
}

// HelloService ...
var HelloService = helloService{
	auth: &Auth{
		appKey: "gyb333",
		appSecret: "20200101",
	},
}

func (h helloService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	if err := h.auth.Check(ctx); err != nil {
		return nil, err
	}

	resp := new(pb.HelloReply)
	resp.Message = "Hello " + in.Name + "."

	return resp, nil
}

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


	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	//tc, err := credentials.NewServerTLSFromFile("conf/server.pem", "conf/server.key")
	tlsServer := gtls.Server{
		CaFile:   "conf/CA/ca.pem",
		CertFile: "conf/CA/serverCA.pem",
		KeyFile:  "conf/CA/server.key",
	}
	tc, err := tlsServer.GetCredentialsByCA()
	if err != nil {
		log.Fatalf("GetTLSCredentialsByCA err: %v", err)
	}



	// 实例化grpc Server
	//s := grpc.NewServer()
	opts := []grpc.ServerOption{
		grpc.Creds(tc),	//添加TLS 证书认证
		//grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		//	grpc_ctxtags.StreamServerInterceptor(),
		//	grpc_opentracing.StreamServerInterceptor(),
		//	grpc_prometheus.StreamServerInterceptor,
		//	grpc_zap.StreamServerInterceptor(zapLogger),
		//	grpc_auth.StreamServerInterceptor(myAuthFunction),
		//	grpc_recovery.StreamServerInterceptor(),
		//)),
		//grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		//	grpc_ctxtags.UnaryServerInterceptor(),
		//	grpc_opentracing.UnaryServerInterceptor(),
		//	grpc_prometheus.UnaryServerInterceptor,
		//	grpc_zap.UnaryServerInterceptor(zapLogger),
		//	grpc_auth.UnaryServerInterceptor(myAuthFunction),
		//	grpc_recovery.UnaryServerInterceptor(),
		//)),
		grpc_middleware.WithUnaryServerChain(
			RecoveryInterceptor,
			LoggingInterceptor,
			otgrpc.OpenTracingServerInterceptor(tracer, otgrpc.LogPayloads()),
		),
	}

	s := grpc.NewServer(opts...)


	// 注册HelloService
	pb.RegisterHelloServer(s, HelloService)

	grpclog.Infoln("Listen on " + Address)

	s.Serve(listen)
}

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("gRPC method: %s, %v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	log.Printf("gRPC method: %s, %v", info.FullMethod, resp)
	return resp, err
}

func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()

	return handler(ctx, req)
}


type Auth struct {
	appKey    string
	appSecret string
}

func (a *Auth) Check(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata.FromIncomingContext err")
	}

	var (
		appKey    string
		appSecret string
	)
	if value, ok := md["app_key"]; ok {
		appKey = value[0]
	}
	if value, ok := md["app_secret"]; ok {
		appSecret = value[0]
	}

	if appKey != a.GetAppKey() || appSecret != a.GetAppSecret() {
		return status.Errorf(codes.Unauthenticated, "invalid token")
	}

	return nil
}

func (a *Auth) GetAppKey() string {
	return a.appKey
}

func (a *Auth) GetAppSecret() string {
	return a.appSecret
}