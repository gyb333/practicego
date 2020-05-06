drop table if exists  t;
CREATE TABLE t (
                     `id` int(11) NOT NULL,
                     `a` int(11) DEFAULT NULL,
                     `b` int(11) DEFAULT NULL,
                     PRIMARY KEY (`id`),
                     KEY `a` (`a`),
                     KEY `b` (`b`)
) ENGINE=InnoDB;

delimiter ;;
create procedure idata()
begin
    declare i int;
    set i=1;
    while(i<=100000)do
            insert into t values(i, i, i);
            set i=i+1;
        end while;
end;;
delimiter ;
call idata();

explain select * from t where a between 10000 and 20000;#使用索引a

#SESSION A                                      SESSION B
start transaction with consistent snapshot ;
                                                delete from t;
                                                call idata();
                                                explain select * from t where a between 10000 and 20000;
commit ;



set long_query_time=0;
select * from t where a between 10000 and 20000; /*Q1*/
select * from t force index(a) where a between 10000 and 20000;/*Q2*/

explain select * from t where a between 10000 and 20000; /*Q1*/
explain select * from t force index(a) where a between 10000 and 20000;/*Q2*/

show index from t;
#对于由于索引统计信息不准确导致的问题，你可以用analyze table来解决。
analyze table t;
#而对于其他优化器误判的情况，你可以在应用端用force index来强行指定索引，也可以通过修改语句来引导优化器

explain select * from t where (a between 1 and 1000) and (b between 50000 and 100000) order by b limit 1;

select * from t where (a between 1 and 1000) and (b between 50000 and 100000) order by b limit 1;
select * from t force index(a) where (a between 1 and 1000) and (b between 50000 and 100000) order by b limit 1;
