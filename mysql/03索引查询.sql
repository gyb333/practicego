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

/*
 "Using index condition":会查找使用了索引，但是需要回表查询数据
 “Using index”：表示的就是使用了覆盖索引；
 "Using where"：在查找使用索引的情况下，需要回表去查询所需的数据
 using index & using where：查找使用了索引，但是需要的数据都在索引列中能找到，所以不需要回表查询数据
 “Using filesort”：表示的就是需要排序，MySQL会给每个线程分配一块内存用于排序，称为sort_buffer。
  Using temporary，表示的是需要使用临时表；
 */

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


create table SUser(
                      ID bigint unsigned primary key,
                      email varchar(64)
)engine=innodb;
alter table SUser add index index1(email);

select count(distinct left(email,4))as L4,
       count(distinct left(email,5))as L5,
       count(distinct left(email,6))as L6,
       count(distinct left(email,7))as L7
from SUser;
#使用前缀索引很可能会损失区分度，需要预先设定一个可以接受的损失比例，比如5%。找出不小于 L * 95%的值，假设这里L6、L7都满足，你就可以选择前缀长度为6。
#使用前缀索引，定义好长度，就可以做到既节省空间，又不用额外增加太多的查询成本。
alter table SUser add index index2(email(6));

#第一种方式是使用倒序存储。如果你存储身份证号的时候把它倒过来存，每次查询的时候，你可以这么写：
#select field_list from t where id_card = reverse('input_id_card_string');
#第二种方式是使用hash字段。你可以在表上再创建一个整数字段，来保存身份证的校验码，同时在这个字段上创建索引。
alter table t add id_card_crc int unsigned, add index(id_card_crc);
#select field_list from t where id_card_crc=crc32('input_id_card') and id_card='input_id_card'
/*
都不支持范围查询。倒序存储的字段上创建的索引是按照倒序字符串的方式排序的，
已经没有办法利用索引方式查出身份证号码在[ID_X, ID_Y]的所有市民了。同样地，hash字段的方式也只能支持等值查询。

区别主要体现在以下三个方面：
    从占用的额外空间来看，倒序存储方式在主键索引上，不会消耗额外的存储空间，而hash字段方法需要增加一个字段。
    当然，倒序存储方式使用4个字节的前缀长度应该是不够的，如果再长一点，这个消耗跟额外这个hash字段也差不多抵消了。

    在CPU消耗方面，倒序方式每次写和读的时候，都需要额外调用一次reverse函数，而hash字段的方式需要额外调用一次crc32()函数。
    如果只从这两个函数的计算复杂度来看的话，reverse函数额外消耗的CPU资源会更小些。

    从查询效率上看，使用hash字段方式的查询性能相对更稳定一些。因为crc32算出来的值虽然有冲突的概率，但是概率非常小，
    可以认为每次查询的平均扫描行数接近1。而倒序存储方式毕竟还是用的前缀索引的方式，也就是说还是会增加扫描行数。
 */