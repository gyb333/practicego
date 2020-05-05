package main

import (
	"github.com/streadway/amqp"
	"os"
	"log"
	. "middleware/MQ/common"
)

/*
go run emit_log.go "kern.critical" "A critical kernal error"
go run emit_log.go "kern.test" "A critical kernal.* error"
go run emit_log.go "test.critical" "A critical *.kernal error"
*/

func main() {
	conn,ch :=GetRabbitConnChan("root","root","Hadoop",5672)
	defer conn.Close()
	defer ch.Close()

	err := ch.ExchangeDeclare("logs_topic","topic",
		true, false,false,false, nil,)
	FailOnError(err, "Failed to declare an exchange")

	body := BodyFrom(os.Args)
	err = ch.Publish("logs_topic",
		SeverityFrom(os.Args),false,false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	FailOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}
