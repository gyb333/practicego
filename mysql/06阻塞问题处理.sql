/*开启MDL的instrument：但是相关instrument并没有开启（MySQL 8.0是默认开启的），其可通过如下两种方式开启，
临时生效
修改performance_schema.setup_instrume nts表，但实例重启后，又会恢复为默认值。
*/
UPDATE performance_schema.setup_instruments SET ENABLED = 'YES', TIMED = 'YES'
WHERE NAME = 'wait/lock/metadata/sql/mdl';
/*
永久生效：在配置文件中设置
[mysqld]
performance-schema-instrument='wait/lock/metadata/sql/mdl=ON'
*/


CREATE TABLE `t` (
                     `id` int(11) NOT NULL,
                     `c` int(11) DEFAULT NULL,
                     PRIMARY KEY (`id`)
) ENGINE=InnoDB;

delimiter ;;
create procedure idata()
begin
    declare i int;
    set i=1;
    while(i<=100000)do
            insert into t values(i,i);
            set i=i+1;
        end while;
end;;
delimiter ;

call idata();

show processlist;
#1.查询长时间不返回
select * from t where id=1;
/*
大概率是表t被锁住了。接下来分析原因的时候，一般都是首先执行一下show processlist命令，看看当前语句处于什么状态。
    等MDL锁：show processlist命令查看Waiting for table metadata lock的示意图。
        有一个线程正在表t上请求或者持有MDL写锁，把select语句堵住了。
session A 通过lock tables命令持有表t的MDL写锁，而session B的查询需要获取MDL读锁。所以，session B进入等待状态。
 */
#Session A                           Session B
lock tables t write;
                            select * from t where id=1;
unlock tables ;

#MySQL启动时需要设置performance_schema=on，相比于设置为off会有10%左右的性能损失
select @@performance_schema;

select * from sys.schema_table_lock_waits;

select object_type,object_schema,object_name,lock_type,lock_duration,lock_status,owner_thread_id
from performance_schema.metadata_locks;



/*
    2.等flush:接下来，我给你举另外一种查询被堵住的情况。
*/
select * from information_schema.processlist where id=1;

#有一个线程正要对表t做flush操作。MySQL里面对表做flush操作的用法，一般有以下两个：
flush tables t with read lock;
flush tables with read lock;
/*
如果指定表t的话，代表的是只关闭表t；如果没有指定具体的表名，则表示关闭MySQL里所有打开的表。
但是正常这两个语句执行起来都很快，除非它们也被别的线程堵住了。
出现Waiting for table flush状态的可能情况是：有一个flush tables命令被别的语句堵住了，然后它又堵住了我们的select语句。
 */
#Session A                      Session B               Session C
select sleep(1) from t;
                            flush tables t;
                                                    select * from t where id=1;
/*
 在session A中，我故意每行都调用一次sleep(1)，这样这个语句默认要执行10万秒，在这期间表t一直是被session A“打开”着。
 然后，session B的flush tables t命令再要去关闭表t，就需要等session A的查询结束。
 这样，session C要再次查询的话，就会被flush 命令堵住了。
 */

/*
等行锁:session A启动了事务，占有写锁，还不提交，是导致session B被堵住的原因。
 */
#Session A                                  Session B
begin ;
update t set c=c+1 where id=1;
                                select * from t where id=1 lock in share mode;

select * from  sys.innodb_lock_waits where locked_table='t' ;
/*
查询慢
    由于字段c上没有索引，这个语句只能走id主键顺序扫描，因此需要扫描5万行。
 */
select * from t where c=50000 limit 1; #坏查询不一定是慢查询。

/*
 带lock in share mode的SQL语句，是当前读，因此会直接读到1000001这个结果，所以速度很快；
 而select * from t where id=1这个语句，是一致性读，因此需要从1000001开始，依次执行undo log，执行了100万次以后，才将1这个结果返回。
 */
#Session A                                              Session B
start transaction with consistent snapshot ;
                                            update t set c=c+1 where id=1;  #执行100万次
select * from t where id=1;
select * from t where id=1 lock in share mode ;#当前读


