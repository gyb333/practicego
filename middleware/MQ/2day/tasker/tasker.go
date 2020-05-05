package main

import (
	"github.com/streadway/amqp"
	"os"
	"log"
	. "middleware/MQ/common"
)

//go run tasker.go  First message.
//go run tasker.go  Second message..
//go run tasker.go  Third message...

func main() {
	conn,ch :=GetRabbitConnChan("root","root","Hadoop",5672)
	defer conn.Close()
	defer ch.Close()


	q, err := ch.QueueDeclare("task_queue",
		true,  false,  false, false, nil,)
	FailOnError(err, "Failed to declare a queue")

	body := BodyFrom(os.Args)
	err = ch.Publish("",q.Name,false, false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
}
