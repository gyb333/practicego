package main

import (
	"./model"
	"./process"
	"./utils"
	"io"
	"log"
	"net"
	"time"
)

func init() {
	//当服务器启动时，我们就去初始化我们的redis的连接池
	var pool=utils.InitPool("docker:6379", 16, 0, 300 * time.Second)
	model.MyUserDao = model.NewUserDao(pool)
}


func main() {
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Println("net.Listen err=", err)
		return
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println("listen.Accept err=", err)
			return
		}
		go processConn(conn)
	}
}

func processConn(conn net.Conn) {
	defer conn.Close()
	process := &process.Process{
		Conn: conn,
	}

	err := process.ProcessData()
	if err != nil {
		if err == io.EOF {
			log.Println("客户端退出，服务器端也退出..")
			return
		}
		log.Println("客户端和服务器通讯协程错误=err", err)
		return
	}
}
