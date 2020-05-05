package common

import (
	"log"
	"fmt"
	"github.com/streadway/amqp"
	"os"
	"strings"
)

func FailOnError(err error, msg string){
	if err != nil {
		log.Fatalf("%s:%s", msg, err)
		panic(fmt.Sprintf("%s:%s", msg, err))
	}
}


func GetRabbitConn() *amqp.Connection{
	RabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%d/", "root", "root", "192.168.56.100", 5672)
	conn, err := amqp.Dial(RabbitUrl)

	FailOnError(err, "Failed to connect to RabbitMQ")
	return conn
}

func GetRabbitConnChan(userName,password,host string,port int) (*amqp.Connection,  *amqp.Channel){
	RabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%d/", userName, password, host, port)
	conn, err := amqp.Dial(RabbitUrl)

	FailOnError(err, "Failed to connect to RabbitMQ")
	ch,err:=conn.Channel()
	FailOnError(err, "Failed to connect to channel")
	return conn,ch
}



func BodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

func SeverityFrom(args []string) string {
	var s string
	if len(args) < 2 || args[1] == "" {
		s = "info"
	}else {
		s = os.Args[1]
	}

	return s
}