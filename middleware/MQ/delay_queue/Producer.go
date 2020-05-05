package main

import (
	"github.com/streadway/amqp"
	. "middleware/MQ/common"
)




func main() {
	conn,ch :=GetRabbitConnChan("root","root","Hadoop",5672)
	defer conn.Close()
	defer ch.Close()

	args:=amqp.Table{
		"x-delay":10000,
	}

	ch.Publish("delayed_exchange","delay.1",true, false,
		amqp.Publishing{
			Headers:args,
			DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body :      []byte("x-delayed-message"),
			//Expiration:"10000",//生产者消息设置超时时间,会阻塞后面的消息
		},)
}
