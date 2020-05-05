package fetcher

import (
	"../../utils"
	"bufio"
	"fmt"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var limiter=time.Tick(100*time.Millisecond)

func Fetch(url string)([]byte,error)  {
	<-limiter
	//resp,err :=http.Get(url)
	// 直接用http.Get(url)进行获取信息会报错：Error: status code 403
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	// 查看自己浏览器中的User-Agent信息（检查元素->Network->User-Agent）
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.80 Safari/537.36")
	resp, err := client.Do(req)
	if resp==nil || resp.Body ==nil{
		err = fmt.Errorf("error:resp.Body %v",resp)
	}
	if err != nil {
		log.Fatalln(err)
		return nil,err
	}
	defer resp.Body.Close()


	if resp.StatusCode!=http.StatusOK{
		return nil,fmt.Errorf("error:status code %d",resp.StatusCode)
	}
	bodyReader :=bufio.NewReader(resp.Body)
	e:=utils.DetermineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}