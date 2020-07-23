package proxy

import (
	"bufio"
	"log"
	"net/http"
	"net/url"
	"time"
)

/*反向代理
一种服务端代理技术，用于隐藏真实服务端节点，并实现负载均衡、缓存、安全校验、协议转换等，常见的有LVS、nginx
 */

// 为了测试，简单的通过当前时间戳取余的方式模拟随机访问后端rs
func GetRandServer()string{
	ports := []string{"8081","8082"}
	n := time.Now().Unix() % 2
	return ports[n]
}

func handler(w http.ResponseWriter,r *http.Request){
	// 解析并修改代理服务
	port := GetRandServer()
	proxyAddr := "http://127.0.0.1:" + port
	proxy ,err := url.Parse(proxyAddr)
	if err != nil{
		log.Println(err)
		return
	}
	r.URL.Scheme = proxy.Scheme
	r.URL.Host = proxy.Host

	// 代理请求
	transport := http.DefaultTransport
	resp ,err := transport.RoundTrip(r)
	if err != nil{
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 将响应结果返回
	for key,value := range resp.Header{
		for _,v := range value{
			w.Header().Add(key,v)
		}
	}
	defer resp.Body.Close()
	bufio.NewReader(resp.Body).WriteTo(w)
}