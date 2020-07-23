package etcd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type ClientDis struct {
	client     *clientv3.Client
	serverList map[string]string
	lock       sync.Mutex
}

func NewClientDis(addr []string, prefixes []string) (cli *ClientDis, err error) {
	cli = &ClientDis{serverList: make(map[string]string)}

	conf := clientv3.Config{Endpoints: addr, DialTimeout: 5 * time.Second}

	cli.client, err = clientv3.New(conf)

	for _, p := range prefixes {
		resp, err := cli.client.Get(context.Background(), p, clientv3.WithPrefix())
		if err == nil && resp != nil && resp.Kvs != nil {
			for i := range resp.Kvs {
				if v := resp.Kvs[i].Value; v != nil {
					cli.SetServiceList(string(resp.Kvs[i].Key), string(resp.Kvs[i].Value))
				}
			}
		}

		go func() {
			rch := cli.client.Watch(context.Background(), p, clientv3.WithPrefix())
			for wResp := range rch {
				for _, ev := range wResp.Events {
					switch ev.Type {
					case mvccpb.PUT:
						cli.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
					case mvccpb.DELETE:
						cli.DelServiceList(string(ev.Kv.Key))
					}
				}
			}
		}()
	}
	return
}

func (this *ClientDis) SetServiceList(key, val string) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.serverList[key] = string(val)
}

func (this *ClientDis) DelServiceList(key string) {
	this.lock.Lock()
	defer this.lock.Unlock()

	delete(this.serverList, key)
}

func (this *ClientDis) ServiceOne(prefix string) string {
	this.lock.Lock()
	defer this.lock.Unlock()

	addrArr := make([]string, 0, 16)
	for k, v := range this.serverList {
		if strings.HasPrefix(k, prefix) {
			addrArr = append(addrArr, v)
		}
	}

	if len(addrArr) == 0 {
		return ""
	}

	// 随机获取一个服务
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return addrArr[r.Intn(len(addrArr))]
}

// 使用方法
func main() {
	// 服务发现
	cli, _ := NewClientDis([]string{"192.168.200.129:2379"}, []string{"/node"})
	fmt.Print(cli.ServiceOne("/node"))
}

