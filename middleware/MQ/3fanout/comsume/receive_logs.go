package main

import (
	"log"
	. "middleware/MQ/common"
)


//将同一个消息发送给多个消费者进行处理
func main() {
	conn,ch :=GetRabbitConnChan("root","root","Hadoop",5672)
	defer conn.Close()
	defer ch.Close()

	err := ch.ExchangeDeclare("logs","fanout",
		true, false, false,false, nil, )
	FailOnError(err, "Failed to declare an exchange")

	//创建临时队列
	q, err := ch.QueueDeclare("",false, false, true, false,nil,)
	FailOnError(err, "Failed to declare a queue")

	//绑定临时队列
	err = ch.QueueBind(q.Name, "","logs", false,nil)
	FailOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(q.Name,"",
		true,   // auto-ack 为true 表示不需要手动发送确认
		false,false,false, nil, )
	FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

//go run receive_logs.go
//go run receive_logs.go