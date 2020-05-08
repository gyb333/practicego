#隐式类型转换
select '10' > 9 ;   #result:1 在MySQL中，字符串和数字做比较的话，是将字符串转换成数字。

CREATE TABLE `tradelog` (
                            `id` int(11) NOT NULL,
                            `tradeid` varchar(32) DEFAULT NULL,
                            `operator` int(11) DEFAULT NULL,
                            `t_modified` datetime DEFAULT NULL,
                            PRIMARY KEY (`id`),
                            KEY `tradeid` (`tradeid`),
                            KEY `t_modified` (`t_modified`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

#条件字段函数操作
select count(*) from tradelog where month(t_modified)=7;
#如果对字段做了函数计算，就用不上索引了，这是MySQL的规定。
#对索引字段做函数操作，可能会破坏索引值的有序性，因此优化器就决定放弃走树搜索功能。
#由于在t_modified字段加了month()函数操作，导致了全索引扫描。为了能够用上索引的快速定位能力
select count(*) from tradelog where
    (t_modified >= '2016-7-1' and t_modified<'2016-8-1') or
    (t_modified >= '2017-7-1' and t_modified<'2017-8-1') or
    (t_modified >= '2018-7-1' and t_modified<'2018-8-1');

#这条语句需要走全表扫描。tradeid的字段类型是varchar(32)，而输入的参数却是整型，所以需要做类型转换。
explain select * from tradelog where tradeid=110717;
explain select * from tradelog where tradeid='110717';
explain select * from tradelog where  CAST(tradeid AS signed int) = 110717;

#隐式字符编码转换
CREATE TABLE `trade_detail` (
                                `id` int(11) NOT NULL,
                                `tradeid` varchar(32) DEFAULT NULL,
                                `trade_step` int(11) DEFAULT NULL, /*操作步骤*/
                                `step_info` varchar(32) DEFAULT NULL, /*步骤信息*/
                                PRIMARY KEY (`id`),
                                KEY `tradeid` (`tradeid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

insert into tradelog values(1, 'aaaaaaaa', 1000, now());
insert into tradelog values(2, 'aaaaaaab', 1000, now());
insert into tradelog values(3, 'aaaaaaac', 1000, now());

insert into trade_detail values(1, 'aaaaaaaa', 1, 'add');
insert into trade_detail values(2, 'aaaaaaaa', 2, 'update');
insert into trade_detail values(3, 'aaaaaaaa', 3, 'commit');
insert into trade_detail values(4, 'aaaaaaab', 1, 'add');
insert into trade_detail values(5, 'aaaaaaab', 2, 'update');
insert into trade_detail values(6, 'aaaaaaab', 3, 'update again');
insert into trade_detail values(7, 'aaaaaaab', 4, 'commit');
insert into trade_detail values(8, 'aaaaaaac', 1, 'add');
insert into trade_detail values(9, 'aaaaaaac', 2, 'update');
insert into trade_detail values(10, 'aaaaaaac', 3, 'update again');
insert into trade_detail values(11, 'aaaaaaac', 4, 'commit');

explain select l.operator from tradelog l , trade_detail d where d.tradeid=l.tradeid and d.id=4;
/*
驱动表trade_detail里id=4的行记为R4，被驱动表tradelog上执行的就是类似这样的SQL 语句：
select operator from tradelog  where traideid =$R4.tradeid.value;
select operator from tradelog  where traideid =CONVERT($R4.tradeid.value USING utf8mb4);
 */

explain select d.* from tradelog l, trade_detail d where d.tradeid=l.tradeid and l.id=2;
/*
第一行显示优化器会先在交易记录表tradelog上查到id=2的行，用上了主键索引，rows=1表示只扫描一行；

第二行key=NULL，表示没有用上交易详情表trade_detail上的tradeid索引，进行了全表扫描。
在这个执行计划里，是从tradelog表中取tradeid字段，再去trade_detail表里查询匹配字段。
因此，我们把tradelog称为驱动表，把trade_detail称为被驱动表，把tradeid称为关联字段。

被驱动表trade_detail类似执行这样的语句：
select * from trade_detail where tradeid=$L2.tradeid.value; ;
$L2.tradeid.value的字符集是utf8mb4。实际上这个语句等同于下面这个写法：
select * from trade_detail  where CONVERT(traideid USING utf8mb4)=$L2.tradeid.value;

连接过程中要求在被驱动表的索引字段上加函数操作，是直接导致对被驱动表做全表扫描的原因。
*/

#把trade_detail表上的tradeid字段的字符集也改成utf8mb4，这样就没有字符集转换的问题了。
alter table trade_detail modify tradeid varchar(32) CHARACTER SET utf8mb4 default null;
#如果数据量比较大， 或者业务上暂时不能做这个DDL的话，那就只能采用修改SQL语句的方法
explain select d.* from tradelog l , trade_detail d where d.tradeid=CONVERT(l.tradeid USING utf8) and l.id=2;



CREATE TABLE `table_a` (
                           `id` int(11) NOT NULL,
                           `b` varchar(10) DEFAULT NULL,
                           PRIMARY KEY (`id`),
                           KEY `b` (`b`)
) ENGINE=InnoDB;
#假设现在表里面，有100万行数据，其中有10万行数据的b的值是’1234567890’

select * from table_a where b='1234567890abcd';

/*这条SQL语句的执行很慢，流程是这样的：

在传给引擎执行的时候，做了字符截断。因为引擎里面这个行只定义了长度是10，所以只截了前10个字节，就是’1234567890’进去做匹配；这样满足条件的数据有10万行；
因为是select *， 所以要做10万次回表；但是每次回表以后查出整行，到server层一判断，b的值都不是’1234567890abcd’;
返回结果是空。
 */


