package tcp

import (
	"encoding/binary"
	"fmt"
	"net"
	"testing"
)

func TestStartServer(t *testing.T) {
	StartServer("localhost:50002")
}

func TestClient(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:50002")
	if err != nil {
		fmt.Println("Error dialing", err.Error())
		return // 终止程序
	}
	var headSize int
	var headBytes = make([]byte, 2)
	s := "hello world"
	content := []byte(s)
	headSize = len(content)
	binary.BigEndian.PutUint16(headBytes, uint16(headSize))
	conn.Write(headBytes)
	conn.Write(content)

	s = "hello go"
	content = []byte(s)
	headSize = len(content)
	binary.BigEndian.PutUint16(headBytes, uint16(headSize))
	conn.Write(headBytes)
	conn.Write(content)

	s = "hello tcp"
	content = []byte(s)
	headSize = len(content)
	binary.BigEndian.PutUint16(headBytes, uint16(headSize))
	conn.Write(headBytes)
	conn.Write(content)
}
