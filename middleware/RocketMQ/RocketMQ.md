收发普通消息:是指RocketMQ 中无特性的消息，区别于有特性的定时消息、顺序消息和事务消息。
顺序消息（FIFO 消息）是RocketMQ 提供的一种严格按照顺序来发布和消费的消息类型。
顺序消息分为两类：
全局顺序：对于指定的一个 Topic，所有消息按照严格的先入先出（First In First Out，简称 FIFO）的顺序进行发布和消费。
分区顺序：对于指定的一个 Topic，所有消息根据 Sharding Key 进行区块分区。 同一个分区内的消息按照严格的 FIFO 顺序进行发布和消费。
Sharding Key 是顺序消息中用来区分不同分区的关键字段，和普通消息的 Key 是完全不同的概念。

Producer 类参数
ProducerModel	设置 Producer 实例化模式，取值说明如下：
CommonProducer：表示普通消息生产者
OrderlyProducer：表示顺序消息生产者
TransProducer：表示事务消息生产者

TransactionStatus	执行本地事务和事务回查的状态，取值说明如下：
CommitTransaction：表示提交事务
RollbackTransaction：表示回滚事务
UnknownTransaction：表示事务状态未知

Consumer 类参数
Model	设置 Consumer 实例的消费模式，取值说明如下：
Clustering：表示集群消费
Broadcasting：表示广播消费

ConsumerModel	设置 Consumer 实例化模式，取值说明如下：
CoCurrently：表示普通消息消费者
Orderly：表示顺序消息消费者

发送事务消息为什么必须要实现回查 Check 机制？
当半事务消息发送完成，但本地事务返回状态为 TransactionStatus.Unknow，或者应用退出导致本地事务未提交任何状态时，从 Broker 的角度看，这条 Half 状态的消息的状态是未知的。因此 Broker 会定期要求发送方能 Check 该 Half 状态消息，并上报其最终状态。

Check 被回调时，业务逻辑都需要做些什么？
事务消息的 Check 方法里面，应该写一些检查事务一致性的逻辑。阿里云 RocketMQ 发送事务消息时需要实现 LocalTransactionChecker 接口，用来处理 Broker 主动发起的本地事务状态回查请求；因此在事务消息的 Check 方法中，需要完成两件事情：

检查该半事务消息对应的本地事务的状态（committed or rollback）。
向 Broker 提交该半事务消息本地事务的状态。