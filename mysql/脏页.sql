/*
把内存里的数据写入磁盘的过程，术语就是flush
当内存数据页跟磁盘数据页内容不一致的时候，我们称这个内存页为“脏页”。
内存数据写入到磁盘后，内存和磁盘上的数据页的内容就一致了，称为“干净页”

平时执行很快的更新操作，其实就是在写内存和日志，而MySQL偶尔“抖”一下的那个瞬间，可能就是在刷脏页（flush）。
那么，什么情况会引发数据库的flush过程呢？
    1.对应的就是InnoDB的redo log写满了。这时候系统会停止所有更新操作，把checkpoint往前推进，redo log留出空间可以继续写。
        把checkpoint位置从CP推进到CP’,就需要将两个点之间的日志,对应的所有脏页都flush到磁盘上。之后从write pos到CP’之间就是可以再写入的redo log的区域。

    2.系统内存不足。当需要新的内存页，而内存不够用的时候，就要淘汰一些数据页，空出内存给别的数据页使用。如果淘汰的是“脏页”，就要先将脏页写到磁盘。
        为什么不能直接把内存淘汰掉，从磁盘读入数据页，然后拿redo log出来应用不就行了？这里其实是从性能考虑的。
        如果刷脏页一定会写盘，就保证了每个数据页有两种状态：一种是内存里存在，内存里就肯定是正确的结果，直接返回；
        另一种是内存里没有数据，就可以肯定数据文件上是正确的结果，读入内存后返回。这样的效率最高。
    3.MySQL认为系统“空闲”的时候。
    4.MySQL正常关闭的情况。MySQL会把内存的脏页都flush到磁盘上，下次MySQL启动的时候，就可以直接从磁盘上读数据，启动速度会很快。
 */

 /*
InnoDB用缓冲池（buffer pool）管理内存，缓冲池中的内存页有三种状态：
    第一种是，还没有使用的；
    第二种是，使用了并且是干净页；
    第三种是，使用了并且是脏页。

当要读入的数据页没有在内存的时候，就必须到缓冲池中申请一个数据页。这时候只能把最久不使用的数据页从内存中淘汰掉：
如果要淘汰的是一个干净页，就直接释放出来复用；但如果是脏页呢，就必须将脏页先刷到磁盘，变成干净页后才能复用。

InnoDB需要有控制脏页比例的机制，来尽量避免这两种情况:
    一个查询要淘汰的脏页个数太多，会导致查询的响应时间明显变长；
    日志写满，更新全部堵住，写性能跌为0，这种情况对敏感业务来说，是不能接受的。

  */
select @@innodb_io_capacity;        #控制刷脏页速度12000
select @@innodb_flush_neighbors;    #值为1的时候会有上述的“连坐”机制，值为0时表示不找邻居，自己刷自己的。
#建议innodb_flush_neighbors的值设置成0。因为这时候IOPS往往不是瓶颈，而“只刷自己”，就能更快地执行完必要的刷脏页操作，减少SQL语句响应时间。

/*
一个内存配置为128GB、innodb_io_capacity设置为20000的大规格实例，正常会建议你将redo log设置成4个1GB的文件。
但如果你在配置的时候不慎将redo log设置成了1个100M的文件，会发生什么情况呢？又为什么会出现这样的情况呢？

每次事务提交都要写redo log，如果设置太小，很快就会被写满，也就是下面这个图的状态，这个“环”将很快被写满，write pos一直追着CP。
这时候系统不得不停止所有更新，去推进checkpoint。你看到的现象就是磁盘压力很小，但是数据库出现间歇性的性能下跌。
 */