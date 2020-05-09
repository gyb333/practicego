/*          加锁规则里面，包含了两个“原则”、两个“优化”和一个“bug”。
原则1：加锁的基本单位是next-key lock。希望你还记得，next-key lock是前开后闭区间。
原则2：查找过程中访问到的对象才会加锁。
优化1：索引上的等值查询，给唯一索引加锁的时候，next-key lock退化为行锁。
优化2：索引上的等值查询，向右遍历时且最后一个值不满足等值条件的时候，next-key lock退化为间隙锁。
一个bug：唯一索引上的范围查询会访问到不满足条件的第一个值为止。

 */
drop table if exists t;
CREATE TABLE `t` (
                     `id` int(11) NOT NULL,
                     `c` int(11) DEFAULT NULL,
                     `d` int(11) DEFAULT NULL,
                     PRIMARY KEY (`id`),
                     KEY `c` (`c`)
) ENGINE=InnoDB;

insert into t values(0,0,0),(5,5,5),(10,10,10),(15,15,15),(20,20,20),(25,25,25);
#案例一：等值查询间隙锁
#Session A                                      Session B                               Session C
begin ;
update t set d=d+1 where id=7;
                                insert into t values (8,8,8);   #blocked
                                                                update t set d=d+1 where id=10;#OK
/*
根据原则1，加锁单位是next-key lock，session A加锁范围就是(5,10]；
同时根据优化2，这是一个等值查询(id=7)，而id=10不满足查询条件，next-key lock退化成间隙锁，因此最终加锁的范围是(5,10)。
session B要往这个间隙里面插入id=8的记录会被锁住，但是session C修改id=10这行是可以的。
 */

#案例二：非唯一索引等值锁
#Session A                                      Session B                               Session C
begin ;
select id from t where c=5 lock in share mode ;
                                                update t set d=d+1 where id=5;#OK
                                                                                insert into t values (7,7,7);#blocked
/*根据原则1，加锁单位是next-key lock，因此会给(0,5]加上next-key lock。
要注意c是普通索引，因此仅访问c=5这一条记录是不能马上停下来的，需要向右遍历，查到c=10才放弃。

根据原则2，访问到的都要加锁，因此要给(5,10]加next-key lock。
但是同时这个符合优化2：等值判断，向右遍历，最后一个值不满足c=5这个等值条件，因此退化成间隙锁(5,10)。

根据原则2 ，只有访问到的对象才会加锁，这个查询使用覆盖索引，并不需要访问主键索引，所以主键索引上没有加任何锁，
这就是为什么session B的update语句可以执行完成。但session C要插入一个(7,7,7)的记录，就会被session A的间隙锁(5,10)锁住。
 */

#案例三：主键索引范围锁
#Session A                                          Session B                               Session C
begin ;
select * from t where id>=10 and id<11 for update ;
                                                    insert into t values (8,8,8);
                                                    insert into t values (13,13,13);#blocked
                                                                                    update t set d=d+1 where id=15;#blocked
/*
开始执行的时候，要找到第一个id=10的行，因此本该是next-key lock(5,10]。
根据优化1，主键id上的等值条件，退化成行锁，只加了id=10这一行的行锁。
范围查找就往后继续找，找到id=15这一行停下来，因此需要加next-key lock(10,15]。
所以，session A这时候锁的范围就是主键索引上，行锁id=10和next-key lock(10,15]。
首次session A定位查找id=10的行的时候，是当做等值查询来判断的，而向右扫描到id=15的时候，用的是范围查询判断。
 */

#案例四：非唯一索引范围锁
#Session A                                          Session B                               Session C
begin ;
select * from t where c>=10 and c<11 for update ;
                                                insert into t values (8,8,8);#blocked
                                                                                update t set d=d+1 where c=15;#blocked
/*
用c=10定位记录的时候，索引c上加了(5,10]这个next-key lock后，由于索引c是非唯一索引，没有优化规则，不会蜕变为行锁，
  因此最终sesion A加的锁是，索引c上的(5,10] 和(10,15] 这两个next-key lock。
  需要扫描到c=15才停止扫描，是合理的，因为InnoDB要扫到c=15，才知道不需要继续往后找了。
 */

#案例五：唯一索引范围锁bug
#Session A                                          Session B                               Session C
begin ;
select * from t where id>10 and id<=15 for update ;
                                                    update t set d=d+1 where id=20;#blocked
                                                                                    insert into t values (16,16,16);#blocked
/*
session A是一个范围查询，按照原则1的话，应该是索引id上只加(10,15]这个next-key lock，并且id是唯一键，所以循环判断到id=15这一行就应该停止了。
但是实现上，InnoDB会往前扫描到第一个不满足条件的行为止，也就是id=20。而且由于这是个范围扫描，因此索引id上的(15,20]这个next-key lock也会被锁上。

session B要更新id=20这一行，是会被锁住的。同样地，session C要插入id=16的一行，也会被锁住。
照理说，这里锁住id=20这一行的行为，其实是没有必要的。因为扫描到id=15，就可以确定不用往后再找了。但实现上还是这么做了，因此我认为这是个bug。
 */

#案例六：非唯一索引上存在"等值"的例子
insert into t values(30,10,30);
/*
虽然有两个c=10，但是它们的主键值id是不同的（分别是10和30），因此这两个c=10的记录之间，也是有间隙的。
delete语句加锁的逻辑，其实跟select ... for update 是类似的
 */
#Session A                                          Session B                               Session C
begin ;
delete from t where c=10;
                            insert into t values (12,12,12);#blocked
                                                                    update t set d=d+1 where c=15;#OK
/*
 session A在遍历的时候，先访问第一个c=10的记录。同样地，根据原则1，这
 里加的是(c=5,id=5)到(c=10,id=10)这个next-key lock。
 然后，session A向右查找，直到碰到(c=15,id=15)这一行，循环才结束。根据优化2，这是一个等值查询，向右查找到了不满足条件的行，
 所以会退化成(c=10,id=10) 到 (c=15,id=15)的间隙锁。
 表示开区间，即(c=5,id=5)和(c=15,id=15)这两行上都没有锁。
 */

#案例七：limit 语句加锁
#Session A                                  Session B
begin ;
delete from t where c=10 limit 2;
                                    insert into t values (12,12,12);#OK

/*session A的delete语句加了 limit 2。表t里c=10的记录其实只有两条，因此加不加limit 2，删除的效果都是一样的，但是加锁的效果却不同。
在遍历到(c=10, id=30)这一行之后，满足条件的语句已经有两条，循环就结束了。
因此，索引c上的加锁范围就变成了从（c=5,id=5)到（c=10,id=30)这个前开后闭区间
在删除数据的时候尽量加limit。这样不仅可以控制删除数据的条数，让操作更安全，还可以减小加锁的范围。
 */

#案例八：一个死锁的例子
#Session A                                                  Session B
begin ;
select id from t where c=10 lock in share mode ;
                                                update t set d=d+1 where c=10;#blocked
insert into t values (8,8,8);#blocked

/*
Session A 启动事务后执行查询语句加lock in share mode，在索引c上加了next-key lock(5,10] 和间隙锁(10,15)；

Session B 的update语句在索引c上加next-key lock(5,10] ,先是加(5,10)的间隙锁，加锁成功；然后加c=10的行锁，被锁住进入等待。
然后session A要再插入(8,8,8)这一行，被session B的间隙锁锁住。由于出现了死锁，InnoDB让session B回滚。
 */

#order by
#Session A                                          Session B                               Sesson C
begin ;
select * from t where c>=15 and c<=20
order by c desc for update ;
                                        insert into t values (11,11,11);#blocked
                                                                                    insert into t values (6,6,6)#blocked
/*
实际上，这里session B和session C的insert 语句都会进入锁等待状态。你可以试着分析一下，出现这种情况的原因是什么？

由于是order by c desc，第一个要定位的是索引c上“最右边的”c=20的行，所以会加上间隙锁(20,25)和next-key lock (15,20]。

在索引c上向左遍历，要扫描到c=10才停下来，所以next-key lock会加到(5,10]，这正是阻塞session B的insert语句的原因。
在扫描过程中，c=20、c=15、c=10这三行都存在值，由于是select *，所以会在主键id上加三个行锁。

因此，session A 的select语句锁的范围就是：
索引c上 (5, 25)；主键索引上id=10、15、20三个行锁。

 */