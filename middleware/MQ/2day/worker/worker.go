package main

import (
	"log"
	"bytes"
	"time"
	. "middleware/MQ/common"
	"os/signal"
	"syscall"
	"os"
)

//go run worker.go
//工作队列，它假设队列中的每一个任务都只会被分发到一个工作者进行处理。
func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	conn,ch :=GetRabbitConnChan("root","root","Hadoop",5672)
	defer conn.Close()
	defer ch.Close()

	q, err := ch.QueueDeclare("task_queue",
		true,false, false, false,nil, )
	FailOnError(err, "Failed to declare a queue")

	//轮询调度、消息持久化和公平分发的特性
	err = ch.Qos(1, 0, false, )// 确保rabbitmq会一个一个发消息
	FailOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(q.Name,
		"",false,false,false,false,nil,)
	FailOnError(err, "Failed to register a consumer")



	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dot_count := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dot_count)
			time.Sleep(t * 10*time.Second)
			log.Printf("Done")
			//手动确认成功
			d.Ack(false)
			//手动确认失败 ,requeue 是否重回队列
			//d.Nack(false,false) //批量执行 变成死信队列 TTL过期也可以变成死信队列
			//d.Reject(false)	//单条执行 	false变成死信队列

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-sig

}
