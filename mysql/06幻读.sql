drop table if exists t;
CREATE TABLE `t` (
                     `id` int(11) NOT NULL,
                     `c` int(11) DEFAULT NULL,
                     `d` int(11) DEFAULT NULL,
                     PRIMARY KEY (`id`),
                     KEY `c` (`c`)
) ENGINE=InnoDB;

insert into t values(0,0,0),(5,5,5), (10,10,10),(15,15,15),(20,20,20),(25,25,25);
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
Select * from t where d=5 for update;   #(5,5,5) 由于d做了全表扫描，导致整个表都被锁定了
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

 跟间隙锁存在冲突关系的，是“往这个间隙中插入一个记录”这个操作。间隙锁之间都不存在冲突关系。
 select ... for share 可多读，但是不可写

select ... for update 仅允许一个读
 */
#Session A                                                          Session B
begin;
select * from t where c=7 lock in share mode ;
                                                    begin;
                                                    select * from t where c=7 for update ;
/*
session B并不会被堵住。因为表t里并没有c=7这个记录，因此session A加的是间隙锁(5,10)。
而session B也是在这个间隙加的间隙锁。它们有共同的目标，即：保护这个间隙，不允许插入值。但，它们之间是不冲突的。

间隙锁和行锁合称next-key lock，每个next-key lock是前开后闭区间。表t初始化以后，(0,0,0),(5,5,5),(10,10,10),(15,15,15),(20,20,20),(25,25,25)
如果用select * from t for update要把整个表所有记录锁起来，就形成了7个next-key lock，
分别是 (-∞,0]、(0,5]、(5,10]、(10,15]、(15,20]、(20, 25]、(25, +suprenum]。
 */
#Session A                                                    Session B
begin ;
select * from t where id=9 for update ;#由于id=9这一行并不存在，会加上间隙锁(5,10);
                                            begin ;
                                            select * from t where id=9 for update ;#由于id=9这一行并不存在，会加上间隙锁(5,10);
                                            insert into t values (9,9,9);#被Session A的间隙锁挡住了，只好进入等待

insert into t values (9,9,9);#被Session B的间隙锁挡住了，只好进入等待
#间隙锁的引入，可能会导致同样的语句锁住更大的范围，这其实是影响了并发度的

