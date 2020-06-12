安装protoc-gen-go
go get -u github.com/golang/protobuf/protoc-gen-go

编译生成protoc-gen-go
cd D:\Go\src\github.com\golang\protobuf\protoc-gen-go
 go build  
 go install  
 
 protoc --version
 编译生成.pb.go文件：
 # 编译hello.proto
 protoc -I . --go_out=plugins=grpc:. ./hello.proto