package main

import (
		. "middleware/MQ/common"
	"github.com/streadway/amqp"
	"os"
	"os/signal"
	"syscall"
	"log"
)


func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	conn,ch :=GetRabbitConnChan("root","root","Hadoop",5672)
	defer conn.Close()
	defer ch.Close()


	//声明delayed Exchange
	exchangeName := "delayed_exchange";
	exchangeType := "x-delayed-message";
	args:=amqp.Table{
		"x-delayed-type":"topic",
	}
	err := ch.ExchangeDeclare(exchangeName,exchangeType,
		true, false,false,false,args)
	FailOnError(err, "Failed to declare an exchange")

	queueName := "delay_queue";
	q, err := ch.QueueDeclare(queueName,
		true,false,false,false, nil)
	FailOnError(err, "Failed to declare a queue")

	routingKey := "delay.#";
	err = ch.QueueBind(q.Name,routingKey,exchangeName,false,nil)
	FailOnError(err, "Failed to bind a queue")

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-sig
}
