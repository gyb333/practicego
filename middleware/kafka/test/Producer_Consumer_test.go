package test

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"log"
	"os"
	"os/signal"
	"sync"
	"testing"
	"time"
)

func TestSyncProducter(t *testing.T)  {
	//var Address = []string{"10.130.138.164:9092","10.130.138.164:9093","10.130.138.164:9094"}

	syncProducer([]string{"hadoop:9092"})
}

//同步消息模式
func syncProducer(address []string)  {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second
	p, err := sarama.NewSyncProducer(address, config)
	if err != nil {
		log.Printf("sarama.NewSyncProducer err, message=%s \n", err)
		return
	}
	defer p.Close()
	topic := "test1"
	srcValue := "sync: this is a message. index=%d"
	for i:=0; i<10; i++ {
		value := fmt.Sprintf(srcValue, i)
		msg := &sarama.ProducerMessage{
			Topic:topic,
			Value:sarama.ByteEncoder(value),
		}
		part, offset, err := p.SendMessage(msg)
		if err != nil {
			log.Printf("send message(%s) err=%s \n", value, err)
		}else {
			fmt.Fprintf(os.Stdout, value + "发送成功，partition=%d, offset=%d \n", part, offset)
		}
		time.Sleep(2*time.Second)
	}
}
func TestASyncProducter(t *testing.T)  {
	SaramaProducer([]string{"hadoop:9092"})
}
func SaramaProducer(address []string)  {
	config := sarama.NewConfig()
	//等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	//随机向partition发送消息
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	//注意，版本设置不对的话，kafka会返回很奇怪的错误，并且无法成功发送消息
	config.Version = sarama.V0_10_0_1

	fmt.Println("start make producer")
	//使用配置,新建一个异步生产者
	producer, e := sarama.NewAsyncProducer(address, config)
	if e != nil {
		fmt.Println(e)
		return
	}
	defer producer.AsyncClose()

	//循环判断哪个通道发送过来数据.
	fmt.Println("start goroutine")
	go func(p sarama.AsyncProducer) {
		for{
			select {
			case  suc:=<-p.Successes():
				fmt.Println("offset: ", suc.Offset, "timestamp: ", suc.Timestamp.String(), "partitions: ", suc.Partition)
			case fail := <-p.Errors():
				fmt.Println("err: ", fail.Err)
			}
		}
	}(producer)

	var value string
	for i:=0;;i++ {
		time.Sleep(500*time.Millisecond)
		time11:=time.Now()
		value = "this is a message 0606 "+time11.Format("15:04:05")
		// 发送的消息,主题。注意：这里的msg必须得是新构建的变量，不然你会发现发送过去的消息内容都是一样的，因为批次发送消息的关系。
		msg := &sarama.ProducerMessage{
			Topic: "0606_test",
		}
		//将字符串转化为字节数组
		msg.Value = sarama.ByteEncoder(value)
		//fmt.Println(value)
		//使用通道发送
		producer.Input() <- msg
	}
}

func TestSaramaConsumer(t *testing.T)  {
	SaramaConsumer([]string{"hadoop:9092"})
}
func SaramaConsumer(address []string)  {
	fmt.Println("start consume")
	config := sarama.NewConfig()

	//提交offset的间隔时间，每秒提交一次给kafka
	config.Consumer.Offsets.CommitInterval = 1 * time.Second

	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	config.Version = sarama.V0_10_0_1

	//consumer新建的时候会新建一个client，这个client归属于这个consumer，并且这个client不能用作其他的consumer
	consumer, err := sarama.NewConsumer(address, config)
	if err != nil {
		panic(err)
	}

	//新建一个client，为了后面offsetManager做准备
	client, err := sarama.NewClient(address, config)
	if err != nil {
		panic("client create error")
	}
	defer client.Close()

	//新建offsetManager，为了能够手动控制offset
	offsetManager,err:=sarama.NewOffsetManagerFromClient("group111",client)
	if err != nil {
		panic("offsetManager create error")
	}
	defer offsetManager.Close()

	//创建一个第2分区的offsetManager，每个partition都维护了自己的offset
	partitionOffsetManager,err:=offsetManager.ManagePartition("0606_test",2)
	if err != nil {
		panic("partitionOffsetManager create error")
	}
	defer partitionOffsetManager.Close()


	fmt.Println("consumer init success")

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	//sarama提供了一些额外的方法，以便我们获取broker那边的情况
	topics,_:=consumer.Topics()
	fmt.Println(topics)
	partitions,_:=consumer.Partitions("0606_test")
	fmt.Println(partitions)

	//第一次的offset从kafka获取(发送OffsetFetchRequest)，之后从本地获取，由MarkOffset()得来
	nextOffset,_:=partitionOffsetManager.NextOffset()
	fmt.Println(nextOffset)

	//创建一个分区consumer，从上次提交的offset开始进行消费
	partitionConsumer, err := consumer.ConsumePartition("0606_test", 2, nextOffset+1)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	fmt.Println("start consume really")

ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed message offset %d\n message:%s", msg.Offset,string(msg.Value))
			//拿到下一个offset
			nextOffset,offsetString:=partitionOffsetManager.NextOffset()
			fmt.Println(nextOffset+1,"...",offsetString)
			//提交offset，默认提交到本地缓存，每秒钟往broker提交一次（可以设置）
			partitionOffsetManager.MarkOffset(nextOffset+1,"modified metadata")

		case <-signals:
			break ConsumerLoop
		}
	}
}


func TestClusterConsumer(t *testing.T)  {
	topic := []string{"test1"}
	var wg = &sync.WaitGroup{}
	wg.Add(2)
	//广播式消费：消费者1
	go clusterConsumer(wg, []string{"hadoop:9092"}, topic, "group-1")
	//广播式消费：消费者2
	go clusterConsumer(wg, []string{"hadoop:9092"}, topic, "group-2")

	wg.Wait()

}
// 支持brokers cluster的消费者
func clusterConsumer(wg *sync.WaitGroup,brokers, topics []string, groupId string)  {
	defer wg.Done()
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// init consumer
	consumer, err := cluster.NewConsumer(brokers, groupId, topics, config)
	if err != nil {
		log.Printf("%s: sarama.NewSyncProducer err, message=%s \n", groupId, err)
		return
	}
	defer consumer.Close()

	// trap SIGINT to trigger a shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Printf("%s:Error: %s\n", groupId, err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("%s:Rebalanced: %+v \n", groupId, ntf)
		}
	}()

	// consume messages, watch signals
	var successes int
Loop:
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				fmt.Fprintf(os.Stdout, "%s:%s/%d/%d\t%s\t%s\n", groupId, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				consumer.MarkOffset(msg, "")  // mark message as processed
				successes++
			}
		case <-signals:
			break Loop
		}
	}
	fmt.Fprintf(os.Stdout, "%s consume %d messages \n", groupId, successes)
}


