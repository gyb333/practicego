安装protoc-gen-go
go get -u github.com/golang/protobuf/protoc-gen-go

#降级protoc-gen-go的版本 
git -C D:/Go/src/github.com/golang/protobuf checkout v1.2.0

编译生成protoc-gen-go
cd D:\Go\src\github.com\golang\protobuf\protoc-gen-go
 go build  
 go install  
 
 protoc --version
 编译生成.pb.go文件：
 # 编译hello.proto
 protoc -I . --go_out=plugins=grpc:. ./hello.proto
 
 将介绍 gRPC 的流式，分为三种类型：
 
 Server-side streaming RPC：服务器端流式 RPC
 Client-side streaming RPC：客户端流式 RPC
 Bidirectional streaming RPC：双向流式 RPC
 
 数据包过大造成的瞬时压力
 接收数据包时，需要所有数据包都接受成功且正确后，才能够回调响应，进行业务处理（无法客户端边发送，服务端边处理）
  protoc -I . --go_out=plugins=grpc:. ./*.proto
  
  
  证书生成
  私钥
  openssl ecparam -genkey -name secp384r1 -out server.key
  自签公钥
  openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650
  
  填写信息
  Country Name (2 letter code) [XX]:
  State or Province Name (full name) []:
  Locality Name (eg, city) [Default City]:
  Organization Name (eg, company) [Default Company Ltd]:
  Organizational Unit Name (eg, section) []:
  Common Name (eg, your name or your server's hostname) []:gyb
  Email Address []:
  
  
  在 gRPC 中，大类可分为两种 RPC 方法，与拦截器的对应关系是：
  
  普通方法：一元拦截器（grpc.UnaryInterceptor）
  流方法：流拦截器（grpc.StreamInterceptor）
  
  基于 CA 的 TLS 证书认证
  生成 Key
  openssl genrsa -out ca.key 2048
  生成密钥
  openssl req -new -x509 -days 7200 -key ca.key -out ca.pem
  填写信息
  Country Name (2 letter code) []:
  State or Province Name (full name) []:
  Locality Name (eg, city) []:
  Organization Name (eg, company) []:
  Organizational Unit Name (eg, section) []:
  Common Name (eg, fully qualified host name) []:gyb333
  Email Address []:
  
  Server
  生成 CSR
  openssl req -new -key server.key -out serverCA.csr
  填写信息
  Country Name (2 letter code) []:
  State or Province Name (full name) []:
  Locality Name (eg, city) []:
  Organization Name (eg, company) []:
  Organizational Unit Name (eg, section) []:
  Common Name (eg, fully qualified host name) []:gyb333
  Email Address []:
  
  Please enter the following 'extra' attributes
  to be sent with your certificate request
  A challenge password []:
  基于 CA 签发
  openssl x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in serverCA.csr -out serverCA.pem
  
  Client
  生成 Key
  openssl ecparam -genkey -name secp384r1 -out client.key
  生成 CSR
  openssl req -new -key client.key -out clientCA.csr
  基于 CA 签发
  openssl x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in clientCA.csr -out clientCA.pem
  
  
  Zipkin 是分布式追踪系统。它的作用是收集解决微服务架构中的延迟问题所需的时序数据。它管理这些数据的收集和查找
  
  docker run -d -p 9411:9411 openzipkin/zipkin
  
  访问 http://127.0.0.1:9411/zipkin/ 检查 Zipkin 是否运行正常