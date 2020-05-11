/*
在两阶段提交的不同时刻，MySQL异常重启会出现什么现象。

写入redo log 处于prepare阶段之后、写binlog之前，发生了崩溃（crash），由于此时binlog还没写，redo log也还没提交，所以崩溃恢复的时候，这个事务会回滚。这时候，binlog还没写，所以也不会传到备库。

如果redo log里面的事务是完整的，也就是已经有了commit标识，则直接提交；
如果redo log里面的事务只有完整的prepare，则判断对应的事务binlog是否存在并完整：
a. 如果是，则提交事务；
b. 否则，回滚事务。

MySQL怎么知道binlog是完整的?
回答：一个事务的binlog是有完整格式的：statement格式的binlog，最后会有COMMIT；row格式的binlog，最后会有一个XID event。
另外，在MySQL 5.6.2版本以后，还引入了binlog-checksum参数，用来验证binlog内容的正确性。
对于binlog日志由于磁盘原因，可能会在日志中间出错的情况，MySQL可以通过校验checksum的结果来发现。

redo log 和 binlog是怎么关联起来的?
回答：它们有一个共同的数据字段，叫XID。崩溃恢复的时候，会按顺序扫描redo log：
如果碰到既有prepare、又有commit的redo log，就直接提交；
如果碰到只有parepare、而没有commit的redo log，就拿着XID去binlog找对应的事务。

处于prepare阶段的redo log加上完整binlog，重启就能恢复，MySQL为什么要这么设计?
回答：数据与备份的一致性有关。在时刻B，也就是binlog写完以后MySQL发生崩溃，这时候binlog已经写入了，之后就会被从库（或者用这个binlog恢复出来的库）使用。
所以，在主库上也要提交这个事务。采用这个策略，主库和备库的数据就保证了一致性。

如果这样的话，为什么还要两阶段提交呢？干脆先redo log写完，再写binlog。崩溃恢复的时候，必须得两个日志都完整才可以。是不是一样的逻辑？
回答：其实，两阶段提交是经典的分布式系统问题，并不是MySQL独有的。
对于InnoDB引擎来说，如果redo log提交完成了，事务就不能回滚（如果这还允许回滚，就可能覆盖掉别的事务的更新）。
而如果redo log直接提交，然后binlog写入的时候失败，InnoDB又回滚不了，数据和binlog日志又不一致了。
两阶段提交就是为了给所有人一个机会，当每个人都说“我ok”的时候，再一起提交。

不引入两个日志，也就没有两阶段提交的必要了。只用binlog来支持崩溃恢复，又能支持归档，不就可以了？
InnoDB引擎使用的是WAL技术，执行事务的时候，写完内存和日志，事务就算完成了。如果之后崩溃，要依赖于日志来恢复数据页。

那能不能反过来，只用redo log，不要binlog？
一个是归档。redo log是循环写，写到末尾是要回到开头继续写的。这样历史日志没法保留，redo log也就起不到归档的作用。
一个就是MySQL系统依赖于binlog。binlog作为MySQL一开始就有的功能，被用在了很多地方。其中，MySQL系统高可用的基础，就是binlog复制。

redo log一般设置多大？
redo log太小的话，会导致很快就被写满，然后不得不强行刷redo log，这样WAL机制的能力就发挥不出来了。
所以，如果是现在常见的几个TB的磁盘的话，就不要太小气了，直接将redo log设置为4个文件、每个文件1GB吧。

正常运行中的实例，数据写入后的最终落盘，是从redo log更新过来的还是从buffer pool更新过来的呢？
    实际上，redo log并没有记录数据页的完整数据，所以它并没有能力自己去更新磁盘数据页，也就不存在“数据最终落盘，是由redo log更新过去”的情况。
如果是正常运行的实例的话，数据页被修改以后，跟磁盘的数据页不一致，称为脏页。最终数据落盘，就是把内存中的数据页写盘。这个过程，甚至与redo log毫无关系。
在崩溃恢复场景中，InnoDB如果判断到一个数据页可能在崩溃恢复的时候丢失了更新，就会将它读到内存，然后让redo log更新内存内容。
更新完成后，内存页变成脏页，就回到了第一种情况的状态。

redo log buffer是什么？是先修改内存，还是先写redo log文件？
begin;
insert into t1 ...
insert into t2 ...
commit;
在执行第一个insert的时候，数据的内存被修改了，redo log buffer也写入了日志。
真正把日志写到redo log文件（文件名是 ib_logfile+数字），是在执行commit语句的时候做的。
*/

CREATE TABLE `like` (
                        `id` int(11) NOT NULL AUTO_INCREMENT,
                        `user_id` int(11) NOT NULL,
                        `liker_id` int(11) NOT NULL,
                        relation_ship smallint not null ,
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `uk_user_id_liker_id` (`user_id`,`liker_id`)
) ENGINE=InnoDB;

CREATE TABLE `friend` (
    id int(11) NOT NULL AUTO_INCREMENT,
  `friend_1_id` int(11) NOT NULL,
  `firned_2_id` int(11) NOT NULL,
  UNIQUE KEY `uk_friend` (`friend_1_id`,`firned_2_id`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;
/*
 在并发场景下，同时有两个人，设置为关注对方，就可能导致无法成功加为朋友关系。
 因为在业务设定里面，这两个逻辑都执行完成以后，是应该在friend表里面插入一行记录的。

 要给“like”表增加一个字段，比如叫作 relation_ship，并设为整型，取值1、2、3。
值是1的时候，表示user_id 关注 liker_id;
值是2的时候，表示liker_id 关注 user_id;
值是3的时候，表示互相关注。

然后，当 A关注B的时候，逻辑改成如下所示的样子：
 比较A和B的大小，如果A<B，则执行下面的逻辑
    begin;
    insert into `like`(user_id, liker_id, relation_ship) values(A, B, 1)
        on duplicate key update relation_ship=relation_ship | 1;
    select relation_ship from `like` where user_id=A and liker_id=B;
    代码中判断返回的 relation_ship，
      如果是1，事务结束，执行 commit
      如果是3，则执行下面这两个语句：
    insert ignore into friend(friend_1_id, friend_2_id) values(A,B);
    commit;

 如果A>B，则执行下面的逻辑
     begin;
    insert into `like`(user_id, liker_id, relation_ship) values(B, A, 2)
        on duplicate key update relation_ship=relation_ship | 2;
    select relation_ship from `like` where user_id=B and liker_id=A;
    代码中判断返回的 relation_ship，
      如果是2，事务结束，执行 commit
      如果是3，则执行下面这两个语句：

    insert ignore into friend(friend_1_id, friend_2_id) values(B,A);
    commit;


 */


CREATE TABLE `t` (
                     `id` int(11) NOT NULL primary key auto_increment,
                     `a` int(11) DEFAULT NULL
) ENGINE=InnoDB;
insert into t values(1,2);
update t set a=2 where id=1;
/*结果显示，匹配(rows matched)了一行，修改(Changed)了0行。
    更新都是先读后写的，MySQL读出数据，发现a的值本来就是2，不更新，直接返回，执行结束；
    MySQL调用了InnoDB引擎提供的“修改为(1,2)”这个接口，但是引擎发现值与原来相同，不更新，直接返回；
    InnoDB认真执行了“把这个值修改成(1,2)"这个操作，该加锁的加锁，该更新的更新。
你觉得实际情况会是以上哪种呢？你可否用构造实验的方式，来证明你的结论？进一步地，可以思考一下，MySQL为什么要选择这种策略呢？

 第一个选项是，MySQL读出数据，发现值与原来相同，不更新，直接返回，执行结束。这里我们可以用一个锁实验来确认。
  SESSION A                                     SESSION B
  begin;
  update t set a=2 where id=1;
                                                update t set a=2 where id=1;#blocked
 session B的update 语句被blocked了，加锁这个动作是InnoDB才能做的，所以排除选项1。

  第二个选项是，MySQL调用了InnoDB引擎提供的接口，但是引擎发现值与原来相同，不更新，直接返回。有没有这种可能呢？
  SESSION A                                     SESSION B
  begin;
  select  * from t where id=1;
                                                update t set a=3 where id=1;
  update t set a=3 where id=1;
  select  * from t where id=1;#(1,3)
  session A的第二个select 语句是一致性读（快照读)，它是不能看见session B的更新的。现在它返回的是(1,3)

  SESSION A                                     SESSION B
  begin;
  select  * from t where id=1;
                                                update t set a=3 where id=1;
  update t set a=3 where id=1 and a=3;
  select  * from t where id=1;#(1,2)

  update t set a=a+1 where id=1 and a=3;
  select  * from t where id=1;#(1,4)
  选项3，即：InnoDB认真执行了“把这个值修改成(1,2)"这个操作，该加锁的加锁，该更新的更新。

 */

