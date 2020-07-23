package proxy

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)
/*正向代理
一种客户端代理技术，用于帮助客户端访问无法直接访问的网络资源，并隐藏客户端IP，常见的场景有***、浏览器HTTP代理
	代理接收客户端请求，复制该请求对象，并根据实际需要配置请求参数
	构造新的请求，发送到服务端，并获取服务端的响应内容
	接收到响应内容后返回给客户端
 */
type proxy struct {}

func (p *proxy)ServeHTTP(w http.ResponseWriter,r *http.Request){
	fmt.Printf("Received request: %s %s %s\n",r.Method,r.Host,r.RemoteAddr)
	transport := http.DefaultTransport

	// 浅拷贝一个request 对象，避免后续修影响了源对象
	req := new(http.Request)
	*req = *r

	// 设置X-Forward-For 头部
	if clientIp,_,err := net.SplitHostPort(r.RemoteAddr);err ==nil{
		if prior,ok := req.Header["X-Forward-For"];ok{
			clientIp = strings.Join(prior,", ") + ", " + clientIp
		}
		req.Header.Set("X-Forward-For",clientIp)
	}

	// 构造新请求
	response,err:=transport.RoundTrip(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 获取响应数据并返回
	for k,v := range response.Header{
		for _,v1 := range v{
			w.Header().Add(k,v1)
		}
	}
	w.WriteHeader(response.StatusCode)
	io.Copy(w,response.Body)
	response.Body.Close()

}

