package httpServer

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"time"
)

type HttpServer struct {
	Host string
	Weight int
}

type LoadBalance struct {
	Index int
	Servers []*HttpServer
}

func NewLoadBalance()*LoadBalance{
	return &LoadBalance{Servers:make([]*HttpServer,0)}
}

func NewHttpServer(host string)*HttpServer{
	return &HttpServer{
		Host:host,
	}
}
func NewHttpServerByWeight(host string,weight int)*HttpServer{
	return &HttpServer{
		Host:host,
		Weight:weight,
	}
}

func (lb *LoadBalance)Add(server *HttpServer){
	lb.Servers = append(lb.Servers,server)
}

// 随机负载均衡
func (lb *LoadBalance)GetHttpServerByRandom()*HttpServer{
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(lb.Servers))
	return lb.Servers[index]
}

/*
加权随机原理：获取到所有节点的权重值，将weight个当前节点Index加到一个[]int，并随机从中获取一个index，例如：
A ： B ： C = 5:2:1 且ABC三个节点的Index分别为0,1,2，那么新建一个如下是切片：
[]int{0,0,0,0,0,1,1,2} ，然后通过rand(len([]int)) 随机拿到一个index
 */
func (lb *LoadBalance)GetHttpServerByRandomWithWeight()*HttpServer{
	var httpServerArr []int
	for index,server := range lb.Servers{
		if server.Weight > 0 {
			for i:=0;i<server.Weight;i++{
				httpServerArr = append(httpServerArr,index)
			}
		}
	}

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(httpServerArr))
	return lb.Servers[httpServerArr[index]]
}
/* 加权随机算法优化版
上面的加权随机算法实现起来比较简单，但存在一个明显弊端，如果weight值的大小将直接影响切片大小，例如5:2 跟 50000:20000 本质上是一样的，
但后者将占用更多的内存空间。因此我们需要对该算法做下优化，将N个节点权重计算出N个区间，然后取随机数rand(weightSum)，
看该数落在哪个区间就返回该区间对应的index值，举个例子：
假设A:B:C = 5:2:1
那么我们先计算出3个区间：5,7(5+2),8(5+2+1)
[0,5) [5,7) [7,8)
然后取rand(5+2+1)，假设获取到的值为6，则落在[5,7) 这个区间，返回index=1 可以看出rand(7)随机数落在各个区间分布如下：
[0,5) ： 0,1,2,3,4
[5,7) ：5,6
[7,8) ：7
正好是5:2:1
 */
// 加权随机优化版
func (lb *LoadBalance)GetHttpServerByRandomWithWeight2()*HttpServer{
	rand.Seed(time.Now().UnixNano())
	// 计算所有节点权重值之和
	weightSum := 0
	for i:=0;i<len(lb.Servers);i++{
		weightSum += lb.Servers[i].Weight
	}
	// 随机数获取
	randNum := rand.Intn(weightSum)

	sum := 0
	for i := 0;i<len(lb.Servers);i++{
		sum += lb.Servers[i].Weight
		// 因为区间是[ ) ，左闭右开，故随机数小于当前权重sum值，则代表落在该区间，返回当前的index
		if randNum < sum {
			return lb.Servers[i]
		}
	}
	return lb.Servers[0]
}


/*
轮询算法
假设有ABC 3台机器，那么请求过来将按照ABCABC 这样的顺顺序将请求反向代理到后端服务器
原理是记录当前的index值，每次请求+1 取模（这里仅演示算法，未考虑线程安全问题，没有加锁）
 */
// 轮询
func (lb *LoadBalance)GetHttpServerByRoundRobin() *HttpServer{
	server := lb.Servers[lb.Index]
	lb.Index = (lb.Index + 1)% len(lb.Servers)
	return server
}

// 加权轮询切片
func (lb *LoadBalance)GetHttpServerByRoundRobinWithWeight() *HttpServer{
	// 加权轮询切片
	indexArr := make([]int,0)
	for index,server := range lb.Servers{
		if server.Weight > 0{
			for i:=0;i<server.Weight;i++{
				indexArr = append(indexArr,index)
			}
		}
	}
	lb.Index = (lb.Index + 1)% len(indexArr)
	fmt.Println(indexArr)
	return lb.Servers[indexArr[lb.Index]]
}

// 加权轮询区间算法
func (lb *LoadBalance)GetHttpServerByRoundRobinWithWeight2()*HttpServer{
	server := lb.Servers[0]
	sum := 0

	for i:=0;i<len(lb.Servers);i++{
		sum += lb.Servers[i].Weight
		if lb.Index < sum{
			server = lb.Servers[i]
			if lb.Index == sum -1 && i != len(lb.Servers)-1{
				lb.Index++
			}else{
				lb.Index = (lb.Index+1) % sum
			}
			fmt.Println(lb.Index)
			break
		}
	}
	return server
}

//ip_hash 算法
// 对客户端IP 做hash 取模得到有一个固定的index，返回固定的httpserver
func (lb *LoadBalance)GetHttpServerByIpHash(ip string) *HttpServer{
	index := int(crc32.ChecksumIEEE([]byte(ip))) % len(lb.Servers)
	return lb.Servers[index]
}

//url_hash 算法
func (lb *LoadBalance) GetHttpServerByUrlHash(url string) *HttpServer{
	index := int(crc32.ChecksumIEEE([]byte(url))) % len(lb.Servers)
	return lb.Servers[index]
}