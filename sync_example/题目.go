package sync_example
// 抽象一个栅栏
type Barrier interface {
	Wait ()
}
// 创建栅栏对象
func NewBarrier (n int) Barrier {

}
// 栅栏的实现类
type barrier struct {

}
// 测试代码
func main () {
	// 创建栅栏对象
	b := NewBarrier(10)
	// 达到的效果：前9个协程调用Wait()阻塞，第10个调用后10个协程全部唤醒
	for i := 0; i < 10; i++ {
		go b.Wait()
	}
}
