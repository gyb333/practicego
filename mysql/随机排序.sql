CREATE TABLE `words` (
                         `id` int(11) NOT NULL AUTO_INCREMENT,
                         `word` varchar(64) DEFAULT NULL,
                         PRIMARY KEY (`id`)
) ENGINE=InnoDB;
truncate table words;
drop procedure if exists iwdata;
delimiter ;;
create procedure iwdata()
begin
    declare i int;
    set i=0;
    while i<300000 do
            insert into words(word) values(concat(char(97+(i%10000 div 1000)), char(97+(i % 1000 div 100)), char(97+(i % 100 div 10)), char(97+(i % 10))));
            set i=i+1;
        end while;
end;;
delimiter ;

call iwdata();
/*
    对于InnoDB表来说，执行全字段排序会减少磁盘访问，因此会被优先选择。
    对于内存表，回表过程只是简单地根据数据行的位置，直接访问内存得到数据，根本不会导致多访问磁盘。
    MySQL这时就会选择rowid排序
 */
select word from words order by rand() limit 3;
explain select word from words order by rand() limit 3;
#显示Using temporary，表示的是需要使用临时表；Using filesort，表示的是需要执行排序操作。
#需要临时表，并且需要在临时表上排序。

/* 这条语句的执行流程是这样的：
    1.创建一个临时表。使用的是memory引擎，表里有两个字段，第一个是double类型，记为字段R，第二个是varchar(64)类型，记为字段W。并且，这个表没有建索引。
    2.从words表中，按主键顺序取出所有的word值。对于每个word值，调用rand()函数生成大于0小于1的随机小数，并把这个随机小数和word分别存入临时表的R和W字段中，到此，扫描行数是10000。
    3.现在临时表有10000行数据了，接下来你要在这个没有索引的内存临时表上，按照字段R排序。
    4.初始化 sort_buffer。sort_buffer中有两个字段，一个是double类型，另一个是整型。
    5.从内存临时表中一行一行地取出R值和位置信息，分别存入sort_buffer中的两个字段里。这个过程要对内存临时表做全表扫描，此时扫描行数增加10000，变成了20000。
    6.在sort_buffer中根据R的值进行排序。注意，这个过程没有涉及到表操作，所以不会增加扫描行数。
    7.排序完成后，取出前三个结果的位置信息，依次到内存临时表中取出word值，返回给客户端。这个过程中，访问了表的三行数据，总扫描行数变成了20003。
 */

# Query_time: 0.215636  Lock_time: 0.000137 Rows_sent: 5  Rows_examined: 600105
select word from words order by rand() limit 5;
#Rows_examined：600105就表示这个语句执行过程中扫描了600105行

/*
如果把一个InnoDB表的主键删掉，是不是就没有主键，就没办法回表了？
其实不是的。如果你创建的表没有主键，或者把一个表的主键删掉了，那么InnoDB会自己生成一个长度为6字节的rowid来作为主键。

rowid:每个引擎用来唯一标识数据行的信息
    对于有主键的InnoDB表来说，这个rowid就是主键ID；
    对于没有主键的InnoDB表来说，这个rowid就是由系统生成的；
    MEMORY引擎不是索引组织表。在这个例子里面，你可以认为它就是一个数组。因此，这个rowid其实就是数组的下标。

order by rand()使用了内存临时表，内存临时表排序的时候使用了rowid排序方法
 */

/*
    是不是所有的临时表都是内存表呢？
    其实不是的。tmp_table_size这个配置限制了内存临时表的大小，默认值是16M。
    如果临时表大小超过了tmp_table_size，那么内存临时表就会转成磁盘临时表。
 */

#为了复现这个过程，我把tmp_table_size设置成1024，把sort_buffer_size设置成 32768, 把 max_length_for_sort_data 设置成16。
select @@tmp_table_size;            #16777216
select @@sort_buffer_size;          #262144
select @@max_length_for_sort_data;  #1024

set tmp_table_size=1024;
set sort_buffer_size=32768;
set max_length_for_sort_data=16;


/* 打开 optimizer_trace，只对本线程有效 */
SET optimizer_trace='enabled=on';

/* 执行语句 */
select word from words order by rand() limit 3;

/* 查看 OPTIMIZER_TRACE 输出 */
SELECT * FROM `information_schema`.`OPTIMIZER_TRACE` G;