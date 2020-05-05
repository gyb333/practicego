package main

import (
	"github.com/streadway/amqp"
		. "middleware/MQ/common"
		"strconv"
	)



func main() {
	conn,ch :=GetRabbitConnChan("root","root","Hadoop",5672)
	defer conn.Close()
	defer ch.Close()

	queueName:="helloQueue"
	q, err := ch.QueueDeclare(queueName,
		true,false, false,  false,nil, )
	FailOnError(err, "Failed to declare q queue")

	body := "Hello"
	for i:=0;i<10;i++{
		err = ch.Publish("", q.Name,   false,false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType: "text/plain",
				Body :      []byte(body+strconv.Itoa(i)),
			},
		)
	}


	FailOnError(err, "Failed to publish a message")

}
