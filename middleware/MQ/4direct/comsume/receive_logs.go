package main

import (
		"os"
	"log"
	. "middleware/MQ/common"
)

//实现了一个可以广播消息给多个接收者的日志系统
//使用direct类型的exchange替代了只能广播消息的fanout类型，让日志系统能够有选择性的接收处理消息
func main(){
	conn :=GetRabbitConn()
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare("logs_direct","direct",
		true,false,false,false,nil,)
	FailOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare("",
		false,false,true,false,nil, )
	FailOnError(err, "Failed to declare a queue")

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [info] [warning] [error]")
		os.Exit(0)
	}

	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s",
			q.Name, "logs_direct", s)
		err = ch.QueueBind(q.Name,s, "logs_direct",false,nil,)
		FailOnError(err, "Failed to bind a queue")
	}

	msgs, err := ch.Consume(q.Name,    "",
		true,false,false,false,nil,)
	FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func(){
		for d:= range msgs{
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To Exit press Ctrl+c")
	<-forever
}


//go run receive_logs.go warning error>logs_from_rabbit.log
//go run receive_logs.go warning error info