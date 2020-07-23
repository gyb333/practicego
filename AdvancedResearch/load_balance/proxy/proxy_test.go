package proxy

import (
	"log"
	"net/http"
	"testing"
)

//修改浏览器代理为 127.0.0.1:9090
/*
$env:http_proxy=http://127.0.0.1:9090
curl -v http://www.baidu.com
 */
func TestProxy(t *testing.T)  {
	server := &http.Server{
		Addr: ":9090",
		Handler: &proxy{},
	}
	if err:=server.ListenAndServe();err != nil{
		log.Fatal("Http proxy server start failed.")
	}
}
