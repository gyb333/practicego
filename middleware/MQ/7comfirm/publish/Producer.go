package main

import (
			. "middleware/MQ/common"
	"github.com/streadway/amqp"
		"log"
	"time"
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

	if err := ch.Confirm(false); err != nil {
		FailOnError(err,"Channel could not be put into confirm mode: %s")
	}
	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	go confirm(confirms)


	exchangeName := "confirm_exchange";
	routingKey := "confirm.qiye";

	ch.Publish(exchangeName,routingKey,false,false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body :      []byte("Send Msg By Confirm ..."),
		},)

	<-sig
}



func confirm(confirms <-chan amqp.Confirmation){
	for{
		ticker := time.NewTicker(100*time.Millisecond)
		select {
		case confirm := <-confirms:
			if confirm.Ack {
				log.Println("confirmed delivery with delivery tag: %d", confirm.DeliveryTag)
			}
		case <- ticker.C:
			log.Println("time out")
		//default:
		//	runtime.Gosched()
		}
	}


}
