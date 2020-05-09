CREATE TABLE `ta` (
                     `id` int(11) NOT NULL,
                     `city` varchar(16) NOT NULL,
                     `name` varchar(16) NOT NULL,
                     `age` int(11) NOT NULL,
                     `addr` varchar(128) DEFAULT NULL,
                     PRIMARY KEY (`id`),
                     KEY `city` (`city`)
) ENGINE=InnoDB;
show variables like '%char%';



select city,name,age from ta where city='杭州' order by name limit 1000  ;
/*全字段排序
1.初始化sort_buffer，确定放入name、city、age这三个字段；
2.从索引city找到第一个满足city='杭州’条件的主键id，也就是图中的ID_X；
3.到主键id索引取出整行，取name、city、age三个字段的值，存入sort_buffer中；
4.从索引city取下一个记录的主键id；
5.重复步骤3、4直到city的值不满足查询条件为止，对应的主键id也就是图中的ID_Y；
6.对sort_buffer中的数据按照字段name做快速排序；
7.按照排序结果取前1000行返回给客户端。
    如果max_length_for_sort_data限制则需要，并按照id的值回到原表中取出city、name和age三个字段返回给客户端。

按name排序”这个动作，可能在内存中完成，也可能需要使用外部排序，这取决于排序所需的内存和参数sort_buffer_size。
sort_buffer_size，就是MySQL为排序开辟的内存（sort_buffer）的大小。
如果要排序的数据量小于sort_buffer_size，排序就在内存中完成。但如果排序数据量太大，内存放不下，则不得不利用磁盘临时文件辅助排序。
  “Using filesort”表示的就是需要排序，MySQL会给每个线程分配一块内存用于排序，称为sort_buffer。
 */

#来确定一个排序语句是否使用了临时文件。
show variables like '%show_compatibility_56%';

set global show_compatibility_56=on;

/* 打开optimizer_trace，只对本线程有效 */
SET optimizer_trace='enabled=on';

/* @a保存Innodb_rows_read的初始值 */
select VARIABLE_VALUE into @a from  information_schema.session_status where variable_name = 'Innodb_rows_read';

/* 执行语句 */
select city, name,age from ta where city='杭州' order by name limit 1000;

/* 查看 OPTIMIZER_TRACE 输出 */
SELECT * FROM `information_schema`.`OPTIMIZER_TRACE` G;

/* @b保存Innodb_rows_read的当前值 */
select VARIABLE_VALUE into @b from information_schema.session_status where variable_name = 'Innodb_rows_read';

/* 计算Innodb_rows_read差值 */
select @b-@a;
/*
examined_rows的值还是4000，表示用于排序的数据是4000行。但是select @b-@a这个语句的值变成5000了。
因为这时候除了排序过程外，在排序完成后，还要根据id去原表取值。由于语句是limit 1000，因此会多读1000行。
 */


/*rowid排序
如果MySQL认为排序的单行长度太大会怎么做呢？
  city、name、age 这三个字段的定义总长度是36，我把max_length_for_sort_data设置为16
 */
SET max_length_for_sort_data = 16;

#创建一个city和name的联合索引，对应的SQL语句是：
alter table ta add index city_user(city, name);
/*整个查询过程的流程就变成了：
1.从索引(city,name)找到第一个满足city='杭州’条件的主键id；
2.到主键id索引取出整行，取name、city、age三个字段的值，作为结果集的一部分直接返回；
3.从索引(city,name)取下一个记录主键id；
4.重复步骤2、3，直到查到第1000条记录，或者是不满足city='杭州’条件时循环结束。
查询过程不需要临时表，也不需要排序
  由于(city,name)这个联合索引本身有序,只需要扫描1000次。
 */

#创建一个city、name和age的联合索引，对应的SQL语句就是：
alter table ta add index city_user_age(city, name, age);
/*整个查询语句的执行流程就变成了：
1.从索引(city,name,age)找到第一个满足city='杭州’条件的记录，取出其中的city、name和age这三个字段的值，作为结果集的一部分直接返回；
2.从索引(city,name,age)取下一个记录，同样取出这三个字段的值，作为结果集的一部分直接返回；
3.重复执行步骤2，直到查到第1000条记录，或者是不满足city='杭州’条件时循环结束。
“Using index”，表示的就是使用了覆盖索引
 */


/*假设你的表里面已经有了city_name(city, name)这个联合索引，然后你要查杭州和苏州两个城市中所有的市民的姓名，并且按名字排序，显示前100条记录。
如果SQL查询语句是这么写的 ：select * from t where city in ('杭州',"苏州") order by name limit 100;
    那么，这个语句执行的时候会有排序过程吗，为什么？
    如果业务端代码由你来开发，需要实现一个在数据库端不需要排序的方案，你会怎么实现呢？
    进一步地，如果有分页需求，要显示第101页，也就是说语句最后要改成 “limit 10000,100”， 你的实现方法又会是什么呢？

虽然(city,name)联合索引，对于单个city内部，name是递增的。但是由于不是要单独地查一个city的值，而是同时查了"杭州"和" 苏州 "两个城市，
因此所有满足条件的name就不是递增的了。也就是说，这条SQL语句需要排序。那怎么避免排序呢？
我们要用到(city,name)联合索引的特性，把这一条语句拆成两条语句，执行流程如下：
    1.执行select * from t where city=“杭州” order by name limit 100; 这个语句是不需要排序的，客户端用一个长度为100的内存数组A保存结果。
    2.执行select * from t where city=“苏州” order by name limit 100; 用相同的方法，假设结果被存进了内存数组B。
    3.现在A和B是两个有序数组，然后你可以用归并排序的思想，得到name最小的前100值，就是我们需要的结果了。

    select id,name from t where city="杭州" order by name limit 10100;
    select id,name from t where city="苏州" order by name limit 10100。
    然后，再用归并排序的方法取得按name顺序第10001~10100的name、id的值，然后拿着这100个id到数据库中去查出所有记录。
 */