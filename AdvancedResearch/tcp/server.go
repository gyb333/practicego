package tcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

/*
TCP的拆包粘包有一般有三种解决方案。

使用定长字节
实际使用中，少于固定字长的，要用字符去填充，空间使用率不够高。

使用分隔符
一般用文本传输的，使用分隔符，IM系统一般对性能要求高，不推荐使用文本传输。

用消息的头字节标识消息内容的长度
可以使用二进制传输，效率高，推荐。下面看看怎么实现
 */
const (
	BYTES_SIZE uint16 = 1024
	HEAD_SIZE  int    = 2
)

func StartServer(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Println("Error listening", err.Error())
		return
	}
	for {
		conn, err := listener.Accept()
		fmt.Println(conn.RemoteAddr())
		if err != nil {
			fmt.Println("Error accepting", err.Error())
			return // 终止程序
		}
		go doConn(conn)
	}
}
//coon->[]byte->Buffer->[]byte  用户态就有三次内存拷贝
func doConn(conn net.Conn) {
	var (
		buffer           = bytes.NewBuffer(make([]byte, 0, BYTES_SIZE))
		bytes            = make([]byte, BYTES_SIZE);
		isHead      bool = true
		contentSize int
		head        = make([]byte, HEAD_SIZE)
		content     = make([]byte, BYTES_SIZE)
	)
	for {
		readLen, err := conn.Read(bytes);
		if err != nil {
			log.Println("Error reading", err.Error())
			return
		}
		_, err = buffer.Write(bytes[0:readLen])
		if err != nil {
			log.Println("Error writing to buffer", err.Error())
			return
		}

		for {
			if isHead {
				if buffer.Len() >= HEAD_SIZE {
					_, err := buffer.Read(head)
					if err != nil {
						fmt.Println("Error reading", err.Error())
						return
					}
					contentSize = int(binary.BigEndian.Uint16(head))
					isHead = false
				} else {
					break
				}
			}
			if !isHead {
				if buffer.Len() >= contentSize {
					_, err := buffer.Read(content[:contentSize])
					if err != nil {
						fmt.Println("Error reading", err.Error())
						return
					}
					fmt.Println(string(content[:contentSize]))
					isHead = true
				} else {
					break
				}
			}
		}
	}
}

//coon->Buffer  用户态只有一次内存拷贝
func doConn2(conn net.Conn) {
	var (
		buffer      = newBuffer(conn, 16)
		headBuf     []byte
		contentSize int
		contentBuf  []byte
	)
	for {
		_, err := buffer.readFromReader()
		if err != nil {
			fmt.Println(err)
			return
		}
		for {
			headBuf, err = buffer.seek(HEAD_SIZE);
			if err != nil {
				break
			}
			contentSize = int(binary.BigEndian.Uint16(headBuf))
			if (buffer.Len() >= contentSize-HEAD_SIZE) {
				contentBuf = buffer.read(HEAD_SIZE, contentSize)
				fmt.Println(string(contentBuf))
				continue
			}
			break
		}
	}
}