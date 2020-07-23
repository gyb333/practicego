package etcd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"log"
	"time"
)

//创建租约注册服务
type ServiceReg struct {
	client        *clientv3.Client
	lease         clientv3.Lease
	leaseResp     *clientv3.LeaseGrantResponse
	cancelFunc    func()
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
}

func NewServiceReg(addr []string, serverList map[string]string, timeout int64) (ser *ServiceReg, err error) {
	ser = &ServiceReg{}

	conf := clientv3.Config{Endpoints: addr, DialTimeout: 5 * time.Second}

	if ser.client, err = clientv3.New(conf); err != nil {
		return
	}

	ser.lease = clientv3.NewLease(ser.client)

	//设置租约时间
	if ser.leaseResp, err = ser.lease.Grant(context.TODO(), timeout); err != nil {
		return
	}

	//设置续租
	var ctx context.Context
	ctx, ser.cancelFunc = context.WithCancel(context.TODO())
	if ser.keepAliveChan, err = ser.lease.KeepAlive(ctx, ser.leaseResp.ID); err != nil {
		return
	}

	// 监听 续租情况
	go func() {
		for {
			select {
			case ch := <-ser.keepAliveChan:
				if ch == nil {
					fmt.Printf("已经关闭续租功能\n")
					return
				} else {
					fmt.Printf("续租成功\n")
				}
			}
		}
	}()

	// 注册服务
	kv := clientv3.NewKV(ser.client)
	for key, val := range serverList {
		_, err = kv.Put(context.TODO(), key, val, clientv3.WithLease(ser.leaseResp.ID))
		if err != nil {
			return
		}
	}

	return
}

// 关闭租约服务
func (this *ServiceReg) Close() {
	this.cancelFunc()
	time.Sleep(2 * time.Second)
	_, _ = this.lease.Revoke(context.TODO(), this.leaseResp.ID)
}

// 使用方法
func main() {
	etcdAddr := []string{"http://192.168.200.129:2379"}
	servAddr := map[string]string{
		"/node/192.168.200.1": "192.168.200.1:50051",
		"/node/192.168.200.2": "192.168.200.2:50051",
	}
	ser, err := NewServiceReg(etcdAddr, servAddr, 10)
	if err != nil {
		log.Fatal(err)
	}
	defer ser.Close()
}

