package example_test

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"log"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

/*
发送事务消息为什么必须要实现回查 Check 机制？
当半事务消息发送完成，但本地事务返回状态为 TransactionStatus.Unknow，或者应用退出导致本地事务未提交任何状态时，从 Broker 的角度看，这条 Half 状态的消息的状态是未知的。因此 Broker 会定期要求发送方能 Check 该 Half 状态消息，并上报其最终状态。

Check 被回调时，业务逻辑都需要做些什么？
事务消息的 Check 方法里面，应该写一些检查事务一致性的逻辑。阿里云 RocketMQ 发送事务消息时需要实现 LocalTransactionChecker 接口，用来处理 Broker 主动发起的本地事务状态回查请求；因此在事务消息的 Check 方法中，需要完成两件事情：

检查该半事务消息对应的本地事务的状态（committed or rollback）。
向 Broker 提交该半事务消息本地事务的状态。
 */

type  TransactionLocalListener struct {
	localTrans       *sync.Map
	transactionIndex int32
}
func NewTransactionListener() *TransactionLocalListener {
	return &TransactionLocalListener{
		localTrans: new(sync.Map),
	}
}

func (dl *TransactionLocalListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	nextIndex := atomic.AddInt32(&dl.transactionIndex, 1)
	log.Printf("nextIndex: %v for transactionID: %v\n", nextIndex, msg.TransactionId)
	status := nextIndex % 3
	dl.localTrans.Store(msg.TransactionId, primitive.LocalTransactionState(status+1))

	log.Printf("dl")
	return primitive.UnknowState
}

func (dl *TransactionLocalListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	log.Printf("%v msg transactionID : %v\n", time.Now(), msg.TransactionId)
	v, existed := dl.localTrans.Load(msg.TransactionId)
	if !existed {
		log.Printf("unknow msg: %v, return Commit", msg)
		return primitive.CommitMessageState
	}
	state := v.(primitive.LocalTransactionState)
	switch state {
	case 1:
		log.Printf("checkLocalTransaction COMMIT_MESSAGE: %v\n", msg)
		return primitive.CommitMessageState
	case 2:
		log.Printf("checkLocalTransaction ROLLBACK_MESSAGE: %v\n", msg)
		return primitive.RollbackMessageState
	case 3:
		log.Printf("checkLocalTransaction unknow: %v\n", msg)
		return primitive.UnknowState
	default:
		log.Printf("checkLocalTransaction default COMMIT_MESSAGE: %v\n", msg)
		return primitive.CommitMessageState
	}
}




func TestTransactionProducer(t *testing.T)  {

	p, _ := rocketmq.NewTransactionProducer(
		NewTransactionListener(),
		producer.WithGroupName(Group),
		producer.WithNameServer([]string{NameServer}),
		producer.WithRetry(1),
	)
	err := p.Start()
	if err != nil {
		log.Printf("start producer error: %s\n", err.Error())
		os.Exit(1)
	}

	for i := 0; i < 10; i++ {
		res, err := p.SendMessageInTransaction(context.Background(),
			primitive.NewMessage("TransactionTopic", []byte("Hello RocketMQ again "+strconv.Itoa(i))))

		if err != nil {
			log.Printf("send message error: %s\n", err)
		} else {
			log.Printf("send message success: result=%s\n", res.String())
		}
	}
	time.Sleep(5 * time.Minute)
	err = p.Shutdown()
	if err != nil {
		log.Printf("shutdown producer error: %s", err.Error())
	}
}


