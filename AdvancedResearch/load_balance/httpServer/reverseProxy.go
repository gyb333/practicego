package httpServer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)
//反向代理代码
type ReveseProxyHandler struct {

}
/*
windows powershell:
for ($counter = 1;$counter -le 10;$counter ++)
{
 curl http://127.0.0.1:8888/reverseproxydemo?id=123
}
liunx:
for i in {0..9};do curl -s http://127.0.0.1:8888/reverseproxydemo?id=123;done
*/
func (rph *ReveseProxyHandler)ServeHTTP(w http.ResponseWriter,r *http.Request){

	//url:=random()
	//url:=randomWithWeight()
	url:=randomWithWeight2()
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(w,r)
}

func random() *url.URL {
	lb := NewLoadBalance()
	lb.Add(NewHttpServer("http://127.0.0.1:8080"))
	lb.Add(NewHttpServer("http://127.0.0.1:8081"))
	lb.Add(NewHttpServer("http://127.0.0.1:8082"))
	lb.Add(NewHttpServer("http://127.0.0.1:8083"))
	lb.Add(NewHttpServer("http://127.0.0.1:8084"))

	url,err := url.Parse(lb.GetHttpServerByRandom().Host)
	if err != nil {
		log.Println("[ERR] url.Parse failed,err:",err)
		return nil
	}
	return url
}

func randomWithWeight()  *url.URL{
	lb := NewLoadBalance()
	lb.Add(NewHttpServerByWeight("http://127.0.0.1:8080",1))
	lb.Add(NewHttpServerByWeight("http://127.0.0.1:8081",2))
	lb.Add(NewHttpServerByWeight("http://127.0.0.1:8082",3))
	lb.Add(NewHttpServerByWeight("http://127.0.0.1:8083",4))
	lb.Add(NewHttpServerByWeight("http://127.0.0.1:8084",5))
	url,err := url.Parse(lb.GetHttpServerByRandomWithWeight().Host)
	if err != nil {
		log.Println("[ERR] url.Parse failed,err:",err)
		return nil
	}
	return url
}
func randomWithWeight2()  *url.URL{
	lb := NewLoadBalance()
	lb.Add(NewHttpServerByWeight("http://127.0.0.1:8080",1))
	lb.Add(NewHttpServerByWeight("http://127.0.0.1:8081",2))
	lb.Add(NewHttpServerByWeight("http://127.0.0.1:8082",3))
	lb.Add(NewHttpServerByWeight("http://127.0.0.1:8083",4))
	lb.Add(NewHttpServerByWeight("http://127.0.0.1:8084",5))
	url,err := url.Parse(lb.GetHttpServerByRandomWithWeight2().Host)
	if err != nil {
		log.Println("[ERR] url.Parse failed,err:",err)
		return nil
	}
	return url
}