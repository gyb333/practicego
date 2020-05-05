package main

import (

	"os"
	"log"

	. "middleware/MQ/common"
	"github.com/streadway/amqp"
)

func main(){
	conn,ch :=GetRabbitConnChan("root","root","Hadoop",5672)
	defer conn.Close()
	defer ch.Close()

	err := ch.ExchangeDeclare("logs_direct","direct",
		true,false,false,false,nil,)
	FailOnError(err, "Failed to declare an exchange")

	body := BodyFrom(os.Args)
	err = ch.Publish("logs_direct", SeverityFrom(os.Args), false,false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	FailOnError(err, "Failed to publish a message")

	log.Printf(" [x] sent %s", body)
}

//go run emit_log.go error "this is a log message"

