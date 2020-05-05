package main
import (
	. "middleware/MQ/common"
	"github.com/streadway/amqp"
	"time"
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

	returns := ch.NotifyReturn(make(chan amqp.Return, 1))
	go handleReturn(returns);

	exchangeName := "return_exchange";
	routingKey := "return.qiye";
	// 发送消息
	msg := "Send Msg And Return Listener By RoutingKey : "

	//mandatory 设置不可路由为true,false 路由不可达消息,自动删除
	ch.Publish(exchangeName,routingKey,true, false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body :      []byte(msg+routingKey),
		},)


	//mandatory设置不可路由为true,因为消费端路由key不正确,导致触发amqp.Return
	routingErrorKey := "error.qiye";
	ch.Publish(exchangeName,routingErrorKey,true, false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body :      []byte(msg+routingErrorKey),
		},)
	<-sig
}

func handleReturn(returns <-chan amqp.Return)  {

	for{
		ticker := time.NewTicker(10*time.Second)
		select {
		case Return := <-returns:
			log.Println(Return)

		case <- ticker.C:
			log.Println("time out")

		case <-time.After(time.Second * 3):
			log.Println("request time out")
		//default:
		//	runtime.Gosched()
		//	log.Println("default")
		}

	}
}