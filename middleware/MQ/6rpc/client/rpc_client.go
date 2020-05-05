package main

import (
	"os"
	"strings"
	"strconv"
	"github.com/streadway/amqp"
	"math/rand"
	"time"
	"log"
	. "middleware/MQ/common"
	"fmt"
)


/*
当Client启动时，会创建一个匿名的、独有的回调队列；

对每一个RPC Request，Client都会设置两个参数：用于标识回调队列的reply_to和用于唯一标识的correlation_id;

Request被发送到rpc_queue队列。

RPC服务器等待rpc_queue的消息，一旦消息到达，处理任务后将响应结果消息发送到reply_to指定的队列；

Client等待callback队列的消息，一旦消息到达，查找与correlation_id匹配的request，然后返回给应用程序。
*/

func randomString(l int) string {
	bytes := make([]byte, l)
	for i:=0; i<l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max - min)
}

func bodyFrom(args []string) int {
	var s string
	if(len(args) < 2 || os.Args[1]==""){
		s = "30"
	}else{
		s = strings.Join(args[1:], " ")
	}

	n, err := strconv.Atoi(s)
	FailOnError(err, "Failed to convert arg to integer")
	return n
}

func fibonacciRPC(n int) (res int, err error) {
	conn,ch :=GetRabbitConnChan("root","root","Hadoop",5672)
	defer conn.Close()
	defer ch.Close()

	q, err := ch.QueueDeclare("",
		false,false,  true, false,nil,)
	FailOnError(err, "Failed to declare a queue")
	msgs , err := ch.Consume(q.Name,
		"",true, false,false,false, nil,)
	FailOnError(err, "Faield to register a consumer")

	corrId := randomString(32)
	fmt.Println(q.Name,corrId)
	err = ch.Publish("","rpc_queue",false,false,
		amqp.Publishing{
			ContentType:        "text/plain",
			CorrelationId:        corrId,
			ReplyTo:            q.Name,
			Body:                []byte(strconv.Itoa(n)),
		})
	FailOnError(err, "Failed to publish a message")

	for d:= range msgs {
		if corrId == d.CorrelationId {
			res, err = strconv.Atoi(string(d.Body))
			FailOnError(err, "Failed to convert body to integer")
			break
		}
	}

	return
}


//go run rpc_client.go 30
func main(){
	rand.Seed(time.Now().UTC().UnixNano())

	n:= bodyFrom(os.Args)

	log.Printf(" [x] Requesting fib(%d)", n)
	res, err := fibonacciRPC(n)
	FailOnError(err, "Failed to handle RPC request")

	log.Printf(" [.] Got %d", res)
}