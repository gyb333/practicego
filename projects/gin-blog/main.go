package main

import (
	"context"
	"fmt"
	_ "gin-blog/docs"
	"gin-blog/models"
	"gin-blog/pkg/gredis"
	"gin-blog/pkg/logging"
	"gin-blog/pkg/setting"
	"gin-blog/routers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Gin-blog API
// @version 1.0
// @description  Golang api of Gin Web
// @termsOfService http://github.com
// @contact.name API Support
// @contact.url http://www.cnblogs.com
// @contact.email 876368840@qq.com
// @host localhost:8000
func main() {
	generatorRun()
	//graceRun()	//liunx优雅的执行
}

/*
不关闭现有连接（正在运行中的程序）
新的进程启动并替代旧进程
新的进程接管新的连接
连接要随时响应用户的请求，当用户仍在请求旧进程时要保持连接，新用户应请求新进程，不可以出现拒绝请求的情况
 */
func graceRun()  {
	//endless.DefaultReadTimeOut = setting.ReadTimeout
	//endless.DefaultWriteTimeOut = setting.WriteTimeout
	//endless.DefaultMaxHeaderBytes = 1 << 20
	//endPoint := fmt.Sprintf(":%d", setting.HTTPPort)
	//
	//server := endless.NewServer(endPoint, routers.InitRouter())
	//server.BeforeBegin = func(add string) {
	//	log.Printf("Actual pid is %d", syscall.Getpid())
	//}
	//
	//err := server.ListenAndServe()
	//if err != nil {
	//	log.Printf("Server err: %v", err)
	//}
	//setting.Setup()
	//models.Setup()
	//logging.Setup()
	//
	//endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeout
	//endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeout
	//endless.DefaultMaxHeaderBytes = 1 << 20
	//endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	//
	//server := endless.NewServer(endPoint, routers.InitRouter())
	//server.BeforeBegin = func(add string) {
	//	log.Printf("Actual pid is %d", syscall.Getpid())
	//}
	//
	//err := server.ListenAndServe()
	//if err != nil {
	//	log.Printf("Server err: %v", err)
	//}
}

func generatorRun()  {
	setting.Setup()
	models.Setup()
	logging.Setup()
	gredis.Setup()

	//router := gin.Default()
	router:=routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),//setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,//setting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,//setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Printf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt,syscall.SIGINT,syscall.SIGTERM)
	<- quit

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}