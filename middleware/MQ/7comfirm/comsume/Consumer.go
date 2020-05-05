package main

import (
	. "middleware/MQ/common"
	"log"
	"os"
	"os/signal"
	"syscall"
)
func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	conn,ch :=GetRabbitConnChan("root","root","Hadoop",5672)
	defer conn.Close()
	defer ch.Close()

	exchangeName := "confirm_exchange";
	exchangeType := "topic";
	err := ch.ExchangeDeclare(exchangeName,exchangeType,
		true,false, false, false,nil,)
	FailOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare("confirm_queue",
		true, false,false, false,nil,)
	FailOnError(err, "Failed to declare a queue")

	routingKey := "confirm.*";
	err = ch.QueueBind(q.Name,routingKey,exchangeName,false,nil)
	FailOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(q.Name,
		"",true,false,false, false,nil,)
	FailOnError(err, "Failed to register a consumer")


	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-sig
}
