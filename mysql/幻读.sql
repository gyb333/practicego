CREATE TABLE `t` (
                     `id` int(11) NOT NULL,
                     `c` int(11) DEFAULT NULL,
                     `d` int(11) DEFAULT NULL,
                     PRIMARY KEY (`id`),
                     KEY `c` (`c`)
) ENGINE=InnoDB;

insert into t values(0,0,0),(5,5,5),
                    (10,10,10),(15,15,15),(20,20,20),(25,25,25);
#这个语句序列是怎么加锁的呢？加的锁又是什么时候释放呢？
    begin;
    select * from t where c=5 for update;
    commit;
#c=5这一行的行锁，还是会等到commit的时候才释放的。
/*
命中d=5这一行，对应的主键id=5，在select 语句执行完成后，id=5这一行会加一个写锁，而且由于两阶段锁协议，这个写锁会在执行commit语句的时候释放。
由于字段d上没有索引，因此这条查询语句会做全表扫描。那么，其他被扫描到的，但是不满足条件的5行记录上，会不会被加锁呢？
 */

#Session A                                          Session B                       Session C
begin;
Select * from t where d=5 for update;   #(5,5,5)
                                        update t set d=5 where id=0;
select * from t where d=5 for update;   #(0,0,5),(5,5,5)
                                                                        insert into t values (1,1,5);
select * from t where d=5 for update;   #(0,0,5),(1,1,5),(5,5,5)
commit;
/*
在可重复读隔离级别下，普通的查询是快照读，是不会看到别的事务插入的数据的。因此，幻读在“当前读”下才会出现。
上面session B的修改结果，被session A之后的select语句用“当前读”看到，不能称为幻读。幻读仅专指“新插入的行”。
 */

#Session A                                          Session B                   Session C
begin;
Select * from t where d=5 for update;   #(5,5,5)
                                        update t set d=5 where id=0;
                                        update t set c=5 where id=0;
select * from t where d=5 for update;   #(0,5,5),(5,5,5)
                                                                        insert into t values (1,1,5);
                                                                        update t set c=5 where id=1;
select * from t where d=5 for update;   #(0,5,5),(1,5,5),(5,5,5)
commit;


#Session A                                          Session B                   Session C
begin;
Select * from t where d=5 for update;   #(5,5,5)
update t set d=100 where d=5;       #(5,5,100)
                                            update t set d=5 where id=0;
                                            update t set c=5 where id=0;
select * from t where d=5 for update;   #(0,5,5),(5,5,100)
                                                                            insert into t values (1,1,5);
                                                                            update t set c=5 where id=1;
select * from t where d=5 for update;   #(0,5,5),(1,5,5),(5,5,5)
commit;
#即使把所有的记录都加上锁，还是阻止不了新插入的记录，这也是为什么“幻读”会被单独拿出来解决的原因。

/*
 如何解决幻读？
产生幻读的原因是，行锁只能锁住行，但是新插入记录这个动作，要更新的是记录之间的“间隙”。
因此，为了解决幻读问题，InnoDB只好引入新的锁，也就是间隙锁(Gap Lock)。
顾名思义，间隙锁，锁的就是两个值之间的空隙。比如文章开头的表t，初始化插入了6个记录，这就产生了7个间隙。
 在一行行扫描的过程中，不仅将给行加上了行锁，还给行两边的空隙，也加上了间隙锁。
 */
