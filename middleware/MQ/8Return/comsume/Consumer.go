package main
import (
	. "middleware/MQ/common"
	"log"
)


func main() {
	conn,ch :=GetRabbitConnChan("root","root","Hadoop",5672)
	defer conn.Close()
	defer ch.Close()

	exchangeName := "return_exchange";
	exchangeType := "topic";
	err := ch.ExchangeDeclare(exchangeName,exchangeType,
		true, false,false,false,nil)
	FailOnError(err, "Failed to declare an exchange")

	queueName := "return_queue";
	q, err := ch.QueueDeclare(queueName,
		true,false,true,false, nil)
	FailOnError(err, "Failed to declare a queue")

	routingKey := "return.*";
	err = ch.QueueBind(q.Name,routingKey,exchangeName,false,nil)
	FailOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(q.Name,
		"",true,false,false, false,nil,)
	FailOnError(err, "Failed to register a consumer")
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
