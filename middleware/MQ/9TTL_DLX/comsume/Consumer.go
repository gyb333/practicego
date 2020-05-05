package main
import (
	. "middleware/MQ/common"
	"log"
	"github.com/streadway/amqp"
)


func main() {
	forever := make(chan bool)
	conn,ch :=GetRabbitConnChan("root","root","Hadoop",5672)
	defer conn.Close()
	defer ch.Close()

	//声明死信队列Exchange和Queue绑定
	ch.ExchangeDeclare("ttl.dlx.exchange","topic",
		true, false,false,false,nil)
	ch.QueueDeclare("ttl.dlx.queue",
		true,false,false,false, nil)
	ch.QueueBind("ttl.dlx.queue","#","ttl.dlx.exchange",false,nil)


	//声明普通Exchange
	exchangeName := "ttl_dlx_exchange";
	exchangeType := "topic";
	err := ch.ExchangeDeclare(exchangeName,exchangeType,
		true, false,false,false,nil)
	FailOnError(err, "Failed to declare an exchange")

	queueName := "ttl_dlx_queue";
	args:=amqp.Table{
		"x-dead-letter-exchange":"ttl.dlx.exchange",
		"x-message-ttl":10000,
	}
	q, err := ch.QueueDeclare(queueName,
		true,false,false,false, args)
	FailOnError(err, "Failed to declare a queue")
	routingKey := "ttl.dlx.*";
	err = ch.QueueBind(q.Name,routingKey,exchangeName,false,nil)
	FailOnError(err, "Failed to bind a queue")



	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
