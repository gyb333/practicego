package main
import (
	. "middleware/MQ/common"
	"github.com/streadway/amqp"
		)

/*
消息变成死信有以下几种情况 :

消息被拒绝(basic.reject/basic.nack) 并且requeue重回队列设置成false
消息TTL过期
队列达到最大长度
死信队列的设置 :

首先要设置死信队列的exchange和queue, 然后进行绑定
Exchange : dlx.exchange
Queue : dlx.queue
RoutingKey : #
然后正常声明交换机, 队列, 绑定, 只不过需要在队列加上一个扩展参数即可 : arguments.put(“x-dead-letter-exchange”, “dlx.exchange”);
这样消息在过期, reject或nack(requeue要设置成false), 队列在达到最大长度时, 消息就可以直接路由到死信队列
 */

func main() {
	conn,ch :=GetRabbitConnChan("root","root","Hadoop",5672)
	defer conn.Close()
	defer ch.Close()


	exchangeName := "ttl_dlx_exchange";
	routingKey := "ttl.dlx.test";
	// 发送消息
	msg := "Send Msg And Return Listener By RoutingKey : "

	//mandatory 设置不可路由为true,false 路由不可达消息,自动删除
	ch.Publish(exchangeName,routingKey,true, false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body :      []byte(msg+routingKey),
			Expiration:"60000",//生产者消息设置超时时间
		},)




}

