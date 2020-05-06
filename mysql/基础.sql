show processlist ;

create table T(ID int primary key, c int);
update T set c=c+1 where ID=2;
/*
在计算机中，磁盘存储数据最小单元是扇区，一个扇区的大小是512字节，
而文件系统（例如XFS/EXT4）的最小单元是块，一个块的大小是4k，
而对于InnoDB存储引擎也有自己的最小储存单元，页（Page），一个页的大小是16K。

如果表T中没有字段k，而你执行了这个语句 select * from T where k=1, 那肯定是会报“不存在这个列”的错误：
“Unknown column ‘k’ in ‘where clause’”。你觉得这个错误是在我们上面提到的哪个阶段报出来的呢？
分析器
 */

/*
当数据库上有多个事务同时执行的时候，就可能出现脏读（dirty read）、不可重复读（non-repeatable read）、幻读（phantom read）的问题，
为了解决这些问题，就有了“隔离级别”的概念。
读未提交是指，一个事务还没提交时，它做的变更就能被别的事务看到。
读提交是指，一个事务提交之后，它做的变更才会被其他事务看到。
可重复读是指，一个事务执行过程中看到的数据，总是跟这个事务在启动时看到的数据是一致的。当然在可重复读隔离级别下，未提交变更对其他事务也是不可见的。
串行化，顾名思义是对于同一行记录，“写”会加“写锁”，“读”会加“读锁”。当出现读写锁冲突的时候，后访问的事务必须等前一个事务执行完成，才能继续执行。
 */
show variables like 'transaction_isolation';
set session transaction isolation level read uncommitted ;
set session transaction isolation level read committed ;
set session transaction isolation level repeatable read ;
set session transaction isolation level serializable ;

create table T(c int) engine=InnoDB;
insert into T(c) values(1);

/*
事务的启动方式
1.显式启动事务语句， begin 或 start transaction。配套的提交语句是commit，回滚语句是rollback。
2.set autocommit=0，将这个线程的自动提交关掉。意味着如果你只执行一个select语句，这个事务就启动了，而且并不会自动提交。
这个事务持续存在直到你主动执行commit 或 rollback 语句，或者断开连接。
建议你总是使用set autocommit=1, 通过显式语句的方式来启动事务。
commit work and chain
 */
show variables like 'autocommit';
select @@autocommit;

select @@completion_type;

#可以在information_schema库的innodb_trx这个表中查询长事务，比如下面这个语句，用于查找持续时间超过60s的事务。
select * from information_schema.innodb_trx where TIME_TO_SEC(timediff(now(),trx_started))>60;

/*
值（completion_type的取值）	描述
NO_CHAIN（或0）	COMMIT并且 ROLLBACK 不受影响。这是默认值。
CHAIN （或1）	commit和rollback之后，会自动开启一个事务。
RELEASE（或2）	commit和rollback之后，会终止当前会话的连接。
 */
create table t(a int, primary key (a))engine=innodb;
set @@completion_type=1;    #CHAIN
select @@completion_type;
begin ;
insert into t select 1;
commit work;        #会自动开启下一个事务
insert into t select 2;
insert into t select 2; #主键冲突，插入失败
rollback;           #会自动开启下一个事务
select * from t;


truncate t;
set @@completion_type=2;#RELEASE
select @@completion_type;
begin ;
insert into t select 3;
commit work ;   #会断开当前连接
select version();
select @@completion_type;#NO_CHAIN

truncate t;
set @@completion_type=1;
select @@completion_type;
begin;
insert into t select 1;
#insert into t select 1;
select * from t;
rollback ;
insert into t select 1;
rollback ;  #CHAIN 可以回滚，NO_CHAIN RELEASE都回滚不了
select * from t;

/*
如何避免长事务对业务的影响？
1.首先，从应用开发端来看：
    确认是否使用了set autocommit=0。这个确认工作可以在测试环境中开展，把MySQL的general_log开起来，然后随便跑一个业务逻辑，通过general_log的日志来确认。
    一般框架如果会设置这个值，也就会提供参数来控制行为，你的目标就是把它改成1。

    确认是否有不必要的只读事务。有些框架会习惯不管什么语句先用begin/commit框起来。
    我见过有些是业务并没有这个需要，但是也把好几个select语句放到了事务中。这种只读事务可以去掉。

    业务连接数据库的时候，根据业务本身的预估，通过SET MAX_EXECUTION_TIME命令，来控制每个语句执行的最长时间，避免单个语句意外执行太长时间。

2.其次，从数据库端来看：
    监控 information_schema.Innodb_trx表，设置长事务阈值，超过就报警/或者kill；

    Percona的pt-kill这个工具不错，推荐使用；

    在业务功能测试阶段要求输出所有的general_log，分析日志行为提前发现问题；

    如果使用的是MySQL 5.6或者更新版本，把innodb_undo_tablespaces设置成2（或更大的值）。
    如果真的出现大事务导致回滚段过大，这样设置后清理起来更方便。
 */