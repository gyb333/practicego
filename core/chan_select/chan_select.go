package chan_select

import "fmt"

/*
chan 为nil 会阻塞发送和接受
无缓冲chan 同步阻塞
缓冲chan 小于缓冲大小，实现异步，当缓冲池满了，同样会阻塞
 */
func Channel()  {
	var nilChan chan struct{}
	//close(nilChan) //panic: close of nil
	nilChan=make(chan struct{})	//初始化
	ichan :=make(chan int,1)
	go func() {
		for{
			fmt.Println(<-ichan)
		}

	}()
	go func() {
		ichan<-1
		ichan<-2
		close(nilChan)
	}()
	<-nilChan
	close(ichan)


}

/*
select:
	如果没有case 会阻塞
	监测chan数据流向
	case 必须为IO操作
	对应异步时间处理，需要在for循环使用
	select 超时处理
	如果多个case都满足条件，会随机选择其中之一来执行
 */
