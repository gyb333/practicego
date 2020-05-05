package main

import (
	"os"
	"log"
	. "middleware/MQ/common"
)


/*
对之前的日志系统进行了改进，
使用direct类型的exchange替代了只能广播消息的fanout类型，让日志系统能够有选择性的接收处理消息。
虽然使用direct类型的exchange提升了日志系统的扩展性，但还是有它的局限性存在，那就是无法配置多重标准的路由。
实现这种灵活性，需要来学习下另外一个功能更综合的topic类型的交换器(exchange).
Topic类型的exchange消息的routing key是有一定限制的，必须是一组使用“.”分开的单词。
单词可以是任意的，但是一般来说以能准确的表达功能的为佳。

Topic exchange很灵活，也很容易用此实现其他类型的功能.
如果将"#"指定为绑定键，那么就会接收所有的消息，相当于fanout类型的广播；
如果通配符"*","#"均不作为绑定键使用，那么其功能实现就等同于direct类型；

接收所有消息：go run receive_logs.go "#"
接收来自"kern"设备的消息：go run receive_logs.go "kern.*"
接收所有以"critical"结尾的消息：go run receive_logs.go "*.critical"
创建多重绑定：go run receive_logs.go "kern.*" "*.critical"
 */


func main() {
	conn,ch :=GetRabbitConnChan("root","root","Hadoop",5672)
	defer conn.Close()
	defer ch.Close()

	err := ch.ExchangeDeclare("logs_topic",
		"topic", true, false,false,false,nil,)
	FailOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare("",    // name
		false, false,true, false,nil,)
	FailOnError(err, "Failed to declare a queue")

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [binding_key]...", os.Args[0])
		os.Exit(0)
	}
	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, "logs_topic", s)
		err = ch.QueueBind(q.Name,   s, "logs_topic", false,nil)
		FailOnError(err, "Failed to bind a queue")
	}

	msgs, err := ch.Consume(q.Name,
		"",true,false,false,false,nil,)
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