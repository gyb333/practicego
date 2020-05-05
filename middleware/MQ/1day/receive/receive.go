package main

import (
		"log"
		. "middleware/MQ/common"
)



func main() {
	conn,ch :=GetRabbitConnChan("root","root","192.168.56.100",5672)
	defer conn.Close()
	defer ch.Close()

	queueName:="helloQueue"
	q, err := ch.QueueDeclare(queueName,
		true,false,false,false, nil,)

	FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(q.Name,
		"", true,false, false, false,nil,)
	FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func(){
		for d:= range msgs{
			log.Printf("Received a message : %s", d.Body)

			//d.Ack(false)		//如果 autoAck 为false 必须手动发送一个确认消息.
		}
	}()

	log.Printf(" [*] Waiting for messages, To exit press CTRL+C")
	<-forever
}


