/*
为什么InnoDB不跟MyISAM一样，也把数字存起来呢？
这是因为即使是在同一个时刻的多个查询，由于多版本并发控制（MVCC）的原因，InnoDB表“应该返回多少行”也是不确定的。

InnoDB是索引组织表，主键索引树的叶子节点是数据，而普通索引树的叶子节点是主键值。所以，普通索引树比主键索引树小很多。
对于count(*)这样的操作，遍历哪个索引树得到的结果逻辑上都是一样的。因此，MySQL优化器会找到最小的那棵树来遍历。
在保证逻辑正确的前提下，尽量减少扫描的数据量，是数据库系统设计的通用法则之一。

count(*)、count(主键id)和count(1) 都表示返回满足条件的结果集的总行数；
而count(字段），则表示返回满足条件的数据行里面，参数“字段”不为NULL的总个数。
*/

/*
count(主键id):InnoDB引擎会遍历整张表，把每一行的id值都取出来，返回给server层。server层拿到id后，判断是不可能为空的，就按行累加。

count(1):InnoDB引擎遍历整张表，但不取值。server层对于返回的每一行，放一个数字“1”进去，判断是不可能为空的，按行累加。
count(1)执行得要比count(主键id)快。因为从引擎返回id会涉及到解析数据行，以及拷贝字段值的操作。

count(字段)：
如果这个“字段”是定义为not null的话，一行行地从记录里面读出这个字段，判断不能为null，按行累加；
如果这个“字段”定义允许为null，那么执行的时候，判断到有可能是null，还要把值取出来再判断一下，不是null才累加。
也就是前面的第一条原则，server层要什么字段，InnoDB就返回什么字段。

count(*)是例外，并不会把全部字段取出来，而是专门做了优化，不取值。count(*)肯定不是null，按行累加。

所以结论是：按照效率排序的话，count(字段)<count(主键id)<count(1)≈count(*)，所以我建议你，尽量使用count(*)。
 */

/*在刚刚讨论的方案中，我们用了事务来确保计数准确。
  由于事务可以保证中间结果不被别的事务读到，因此修改计数值和插入新记录的顺序是不影响逻辑结果的。
  但是，从并发系统性能的角度考虑，你觉得在这个事务序列里，应该先插入操作记录，还是应该先更新计数表呢？

用一个计数表记录一个业务表的总行数，在往业务表插入数据的时候，需要给计数值加1。逻辑实现上是启动一个事务，执行两个语句：
    insert into 数据表；
    update 计数表，计数值加1。
    从系统并发能力的角度考虑，怎么安排这两个语句的顺序。

 */
CREATE TABLE `rows_stat` (
                             `table_name` varchar(64) NOT NULL,
                             `row_count` int(10) unsigned NOT NULL,
                             PRIMARY KEY (`table_name`)
) ENGINE=InnoDB;
/*
在更新计数表的时候，一定会传入where table_name=$table_name，使用主键索引，更新加行锁只会锁在一行上。
而在不同业务表插入数据，是更新不同的行，不会有行锁。

把update计数表放后面，如果把update计数表放到事务的第一个语句，同时插入数据的话，等待时间会更长。
 */