package main

import (
	"github.com/streadway/amqp"
	"strconv"
	"log"
	. "middleware/MQ/common"
)



//go run rpc_server.go

func fib(n int) int {
	if n== 0 {
		return 0
	}else if n==1 {
		return 1
	}else {
		return fib(n-1) + fib(n-2)
	}
}
/*
RPC server对Client请求的响应同样需要通过消息队列来传递，可以对每一次请求创建一个回调队列，
但这种方式效率很低，更好的方式是：对于每一个客户端只创建一个回调队列。
但这样会带来一个问题：回调队列接收到一个response之后，如何确定其对应的request？这就需要 correlataion_id来标识。客户端在request中添加一个唯一的correlation_id，在接收到服务器返回的response时，根据该值来确定与之匹配的request并处理。如果未能找到与之匹配的correlation_id，说明该response并不属于当前client的请求，为了安全起见，将其忽略即可。

我们可能会问：为什么在没有找到与之匹配的correlation_id时是将其忽略而不是失败报错？
这是考虑到服务端的竞争条件：假设RPC server在发送response后宕机了，而此时却没能对当前request发出确认消息(ack).
如果这种情况出现，该请求还在队列中会被再次派发。因此当前Request会在服务端处理两次，也会给客户端发送两次Response，
故而，client要能处理重复的response，而server端对于Request需要实现幂等。
*/

func main(){
	conn,ch :=GetRabbitConnChan("root","root","Hadoop",5672)
	defer conn.Close()
	defer ch.Close()

	q, err := ch.QueueDeclare("rpc_queue",
		false,false,false, false,nil, )
	FailOnError(err, "Failed to declare a queue")

	err = ch.Qos(1,0,false,)
	FailOnError(err, "Failed to set Qos")

	msgs, err := ch.Consume(q.Name, "",
		false, false,false, false,nil, )
	FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			n, err := strconv.Atoi(string(d.Body))
			FailOnError(err, "Failed to convert body to an integer")

			log.Printf(" [.] fib(%d)", n)
			response := fib(n)

			err = ch.Publish("",
				d.ReplyTo, false, false,
				amqp.Publishing{
					ContentType :    "text/plain",
					CorrelationId:    d.CorrelationId,
					Body:            []byte(strconv.Itoa(response)),

				})
			FailOnError(err, "Failed to publish a message")

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Awaiting RPC reqeusts")
	<-forever
}