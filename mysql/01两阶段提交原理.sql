/*
 binlog的写入机制：事务执行过程中，先把日志写到binlog cache，事务提交的时候，再把binlog cache写到binlog文件中。

一个事务的binlog是不能被拆开的，因此不论这个事务多大，也要确保一次性写入。这就涉及到了binlog cache的保存问题。

系统给binlog cache分配了一片内存，每个线程一个，参数 binlog_cache_size用于控制单个线程内binlog cache所占内存的大小。如果超过了这个参数规定的大小，就要暂存到磁盘。

事务提交的时候，执行器把binlog cache里的完整事务写入到binlog中，并清空binlog cache。

 每个线程有自己binlog cache，但是共用同一份binlog文件。

write:指把日志写入到文件系统的page cache，并没有把数据持久化到磁盘，所以速度比较快。
fsync:将数据持久化到磁盘的操作。一般情况下，我们认为fsync才占磁盘的IOPS。

write 和fsync的时机，是由参数sync_binlog控制的：
    sync_binlog=0的时候，表示每次提交事务都只write，不fsync；
    sync_binlog=1的时候，表示每次提交事务都会执行fsync；
    sync_binlog=N(N>1)的时候，表示每次提交事务都write，但累积N个事务后才fsync。

因此，在出现IO瓶颈的场景里，将sync_binlog设置成一个比较大的值，可以提升性能。在实际的业务场景中，考虑到丢失日志量的可控性，
一般不建议将这个参数设成0，比较常见的是将其设置为100~1000中的某个数值。
但是，将sync_binlog设置为N，对应的风险是：如果主机发生异常重启，会丢失最近N个事务的binlog日志。
 */

/*
如果你想提升binlog组提交的效果，可以通过设置 binlog_group_commit_sync_delay 和 binlog_group_commit_sync_no_delay_count来实现。

binlog_group_commit_sync_delay参数，表示延迟多少微秒后才调用fsync;

binlog_group_commit_sync_no_delay_count参数，表示累积多少次以后才调用fsync。

 */

/*
redo log的写入机制:事务在执行过程中，生成的redo log是要先写到redo log buffer的。
    redo log buffer里面的内容，是不是每次生成后都要直接持久化到磁盘呢？不需要。
    如果事务执行期间MySQL发生异常重启，那这部分日志就丢了。由于事务并没有提交，所以这时日志丢了也不会有损失。

    事务还没提交的时候，redo log buffer中的部分日志有没有可能被持久化到磁盘呢？确实会有。

redo log可能存在的三种状态
    1.存在redo log buffer中，物理上是在MySQL进程内存中
    2.写到磁盘(write)，但是没有持久化（fsync)，物理上是在文件系统的page cache里面，也就是图中的黄色部分；
    3.持久化到磁盘，对应的是hard disk
日志写到redo log buffer是很快的，wirte到page cache也差不多，但是持久化到磁盘的速度就慢多了。

为了控制redo log的写入策略，InnoDB提供了innodb_flush_log_at_trx_commit参数，它有三种可能取值：
    1.设置为0的时候，表示每次事务提交时都只是把redo log留在redo log buffer中;
    2.设置为1的时候，表示每次事务提交时都将redo log直接持久化到磁盘；
    3.设置为2的时候，表示每次事务提交时都只是把redo log写到page cache。

1.InnoDB有一个后台线程，每隔1秒，就会把redo log buffer中的日志，调用write写到文件系统的page cache，然后调用fsync持久化到磁盘。
注意，事务执行中间过程的redo log也是直接写在redo log buffer中的，这些redo log也会被后台线程一起持久化到磁盘。
也就是说，一个没有提交的事务的redo log，也是可能已经持久化到磁盘的。

2.redo log buffer占用的空间即将达到 innodb_log_buffer_size一半的时候，后台线程会主动写盘。
注意，由于这个事务并没有提交，所以这个写盘动作只是write，而没有调用fsync，也就是只留在了文件系统的page cache。

3.并行的事务提交的时候，顺带将这个事务的redo log buffer持久化到磁盘。假设一个事务A执行到一半，已经写了一些redo log到buffer中，
这时候有另外一个线程的事务B提交，如果innodb_flush_log_at_trx_commit设置的是1，那么按照这个参数的逻辑，
事务B要把redo log buffer里的日志全部持久化到磁盘。就会带上事务A在redo log buffer里的日志一起持久化到磁盘。

innodb_flush_log_at_trx_commit设置成1，那么redo log在prepare阶段就要持久化一次，
因为有一个崩溃恢复逻辑是要依赖于prepare 的redo log，再加上binlog来恢复的

每秒一次后台轮询刷盘，再加上崩溃恢复这个逻辑，InnoDB就认为redo log在commit的时候就不需要fsync了，只会write到文件系统的page cache中就够了。

通常我们说MySQL的“双1”配置，指的就是sync_binlog和innodb_flush_log_at_trx_commit都设置成 1。
一个事务完整提交前，需要等待两次刷盘，一次是redo log（prepare 阶段），一次是binlog。

 */

/*WAL机制主要得益于两个方面：

redo log 和 binlog都是顺序写，磁盘的顺序写比随机写速度要快；

组提交机制，可以大幅度降低磁盘的IOPS消耗。
 */

/*如果你的MySQL现在出现了性能瓶颈，而且瓶颈在IO上，可以通过哪些方法来提升性能呢？针对这个问题，可以考虑以下三种方法：

设置 binlog_group_commit_sync_delay 和 binlog_group_commit_sync_no_delay_count参数，减少binlog的写盘次数。
这个方法是基于“额外的故意等待”来实现的，因此可能会增加语句的响应时间，但没有丢失数据的风险。

将sync_binlog 设置为大于1的值（比较常见是100~1000）。这样做的风险是，主机掉电时会丢binlog日志。

将innodb_flush_log_at_trx_commit设置为2。这样做的风险是，主机掉电的时候会丢数据。

我不建议你把innodb_flush_log_at_trx_commit 设置成0。表示redo log只保存在内存中，这样MySQL本身异常重启也会丢数据，风险太大。
而redo log写到文件系统的page cache的速度也是很快的，设置成2跟设置成0其实性能差不多，但这样做MySQL异常重启时就不会丢数据了，相比之下风险会更小。

 */

/*
问题1：执行一个update语句以后，我再去执行hexdump命令直接查看ibd文件内容，为什么没有看到数据有改变呢？
回答：这可能是因为WAL机制的原因。update语句执行完成后，InnoDB只保证写完了redo log、内存，可能还没来得及将数据写到磁盘。

问题2：为什么binlog cache是每个线程自己维护的，而redo log buffer是全局共用的？
回答：MySQL这么设计的主要原因是，binlog是不能“被打断的”。一个事务的binlog必须连续写，因此要整个事务完成后，再一起写到文件里。
而redo log并没有这个要求，中间有生成的日志可以写到redo log buffer中。redo log buffer中的内容还能“搭便车”，其他事务提交的时候可以被一起写到磁盘中。

问题3：事务执行期间，还没到提交阶段，如果发生crash的话，redo log肯定丢了，这会不会导致主备不一致呢？
回答：不会。因为这时候binlog 也还在binlog cache里，没发给备库。crash以后redo log和binlog都没有了，从业务角度看这个事务也没有提交，所以数据是一致的。

问题4：如果binlog写完盘以后发生crash，这时候还没给客户端答复就重启了。等客户端再重连进来，发现事务已经提交成功了，这是不是bug？
回答：不是。
你可以设想一下更极端的情况，整个事务都提交成功了，redo log commit完成了，备库也收到binlog并执行了。但是主库和客户端网络断开了，
导致事务成功的包返回不回去，这时候客户端也会收到“网络断开”的异常。这种也只能算是事务成功的，不能认为是bug。

实际上数据库的crash-safe保证的是：
    如果客户端收到事务成功的消息，事务就一定持久化了；
    如果客户端收到事务失败（比如主键冲突、回滚等）的消息，事务就一定失败了；
    如果客户端收到“执行异常”的消息，应用需要重连后通过查询当前状态来继续后续的逻辑。此时数据库只需要保证内部（数据和日志之间，主库和备库之间）一致就可以了。

如果一个数据库是被客户端的压力打满导致无法响应的，重启数据库是没用的。
这个问题是因为重启之后，业务请求还会再发。而且由于是重启，buffer pool被清空，可能会导致语句执行得更慢。

有时候一个表上会出现多个单字段索引，这样就可能出现优化器选择索引合并算法的现象。
但实际上，索引合并算法的效率并不好。而通过将其中的一个索引改成联合索引的方法，是一个很好的应对方案。

 */

/*你的生产库设置的是“双1”吗？ 如果平时是的话，你有在什么场景下改成过“非双1”吗？你的这个操作又是基于什么决定的？
另外，我们都知道这些设置可能有损，如果发生了异常，你的止损方案是什么？

业务高峰期。一般如果有预知的高峰期，DBA会有预案，把主库设置成“非双1”。
备库延迟，为了让备库尽快赶上主库。
用备份恢复主库的副本，应用binlog的过程，这个跟上一种场景类似。
批量导入数据的时候。
一般情况下，把生产库改成“非双1”配置，是设置innodb_flush_logs_at_trx_commit=2、sync_binlog=1000。

 */