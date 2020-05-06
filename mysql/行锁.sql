/*
在InnoDB事务中，行锁是在需要的时候才加上的，但并不是不需要了就立刻释放，而是要等到事务结束时才释放。这个就是两阶段锁协议。
如果你的事务中需要锁多个行，要把最可能造成锁冲突、最可能影响并发度的锁尽量往后放。
在InnoDB中，innodb_lock_wait_timeout的默认值是50s，意味着如果采用第一个策略，当出现死锁以后，第一个被锁住的线程要过50s才会超时退出，然后其他线程才有可能继续执行。对于在线服务来说，这个等待时间往往是无法接受的。
 */

 /*
  如果你要删除一个表里面的前10000行数据，有以下三种方法可以做到：
    第一种，直接执行delete from T limit 10000;
    第二种，在一个连接中循环执行20次 delete from T limit 500;
    第三种，在20个连接中同时执行delete from T limit 500。
你会选择哪一种方法呢？为什么呢？
第一种方式（即：直接执行delete from T limit 10000）里面，单个语句占用时间长，锁的时间也比较长；而且大事务还会导致主从延迟。
第三种方式（即：在20个连接中同时执行delete from T limit 500），会人为造成锁冲突。
  */

/*
 InnoDB里面每个事务有一个唯一的事务ID，叫作transaction id。它是在事务开始的时候向InnoDB的事务系统申请的，是按申请顺序严格递增的。
而每行数据也都是有多个版本的。每次事务更新数据的时候，都会生成一个新的数据版本，并且把transaction id赋值给这个数据版本的事务ID，记为row trx_id。
同时，旧的数据版本要保留，并且在新的数据版本中，能够有信息可以直接拿到它。也就是说，数据表中的一行记录，其实可能有多个版本(row)，每个版本有自己的row trx_id。
 */
CREATE TABLE `t` (
                     `id` int(11) NOT NULL,
                     `k` int(11) DEFAULT NULL,
                     PRIMARY KEY (`id`)
) ENGINE=InnoDB;
insert into t(id, k) values(1,1),(2,2);
#事务A                                            事务B                                         事务C
start transaction with consistent snapshot ;
                                            start transaction with consistent snapshot ;
                                                                                        update t set k=k+1 where id=1;
                                            update t set k=k+1 where id=1;
                                            select k from t where id=1;#3
select k from t where id=1;#1
commit ;
                                            commit ;
/*一致性读
对于当前事务的启动瞬间来说，一个数据版本的row trx_id，有以下几种可能：
如果落在绿色部分，表示这个版本是已提交的事务或者是当前事务自己生成的，这个数据是可见的；
如果落在红色部分，表示这个版本是由将来启动的事务生成的，是肯定不可见的；
如果落在黄色部分，那就包括两种情况
a. 若 row trx_id在数组中，表示这个版本是由还没提交的事务生成的，不可见；
b. 若 row trx_id不在数组中，表示这个版本是已经提交了的事务生成的，可见。

  对于一个事务视图来说，除了自己的更新总是可见以外，有三种情况：
    版本未提交，不可见；
    版本已提交，但是是在视图创建后提交的，不可见；
    版本已提交，而且是在视图创建前提交的，可见。
 */

 /*当前读
事务B的视图数组是先生成的，之后事务C才提交，不是应该看不见(1,2)吗，怎么能算出(1,3)来？
   更新数据都是先读后写的，而这个读，只能读当前的值，称为“当前读”（current read）。除了update语句外，select语句如果加锁，也是当前读。
  */
select k from t where id=1 lock in share mode;#当前读 读锁（S锁，共享锁）
select k from t where id=1 for update;#当前读  写锁（X锁，排他锁）。
#事务A                                            事务B                                         事务C
start transaction with consistent snapshot ;
                                            start transaction with consistent snapshot ;
                                                                                        start transaction with consistent snapshot ;
                                                                                        update t set k=k+1 where id=1;#2
                                            update t set k=k+1 where id=1;#阻塞 当前读
                                            select k from t where id=1;#3
                                                                                        commit ;
select k from t where id=1;#1
commit ;
                                            commit ;
/*事务C’的不同是，更新后并没有马上提交，在它提交前，事务B的更新语句先发起了。
  虽然事务C’还没提交，但是(1,2)这个版本也已经生成了，并且是当前的最新版本。那么，事务B的更新语句会怎么处理呢？
  事务C’没提交，也就是说(1,2)这个版本上的写锁还没释放。而事务B是当前读，必须要读最新版本，而且必须加锁，
  因此就被锁住了，必须等到事务C’释放这个锁，才能继续它的当前读。
 */
 /*
  读提交的逻辑和可重复读的逻辑类似，它们最主要的区别是：
在可重复读隔离级别下，只需要在事务开始的时候创建一致性视图，之后事务里的其他查询都共用这个一致性视图；
在读提交隔离级别下，每一个语句执行前都会重新算出一个新的视图。
  */

#事务隔离级别是可重复读。现在，我要把所有“字段c和id值相等的行”的c值清零，但是却发现了一个“诡异”的、改不掉的情况。请你构造出这种情况，并说明其原理。
CREATE TABLE td (
                            `id` int(11) NOT NULL,
                            `c` int(11) DEFAULT NULL,
                            PRIMARY KEY (`id`)
       ) ENGINE=InnoDB;
truncate td;
insert into td(id, c) values(1,1),(2,2),(3,3),(4,4);
#数据无法修改 的场景
#SESSION A                                                  #SESSION B                      #Session b
                                                                                            begin;
begin ;
select * from td;
#select sleep(10);
                                                        update td set c=c+1;                update td set c=c+1;
update td set c=0 where id =c;
                                                                                            commit ;#任意位置都行
select * from td;

