package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func CheckError(err error) {
	if err != nil{
		log.Fatalln(err)
	}
}


func TestNillPoint(t *testing.T)  {
	resp, err := http.Get("https://api.ipify.org?format=json")
	// 关闭 resp.Body 的正确姿势
	if resp != nil {
		defer resp.Body.Close()
	}
	CheckError(err)
	defer resp.Body.Close() // 绝大多数情况下的正确关闭方式
	body, err := ioutil.ReadAll(resp.Body)
	CheckError(err)
	//_, err = io.Copy(ioutil.Discard, resp.Body) // 手动丢弃读取完毕的数据
	fmt.Println(string(body))
}

func TestKeepAlive(t *testing.T)  {
	req, err := http.NewRequest("GET", "http://golang.org", nil)
	CheckError(err)

	req.Close = true
	//req.Header.Add("Connection", "close") // 等效的关闭方式
	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	CheckError(err)
	fmt.Println(resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	CheckError(err)
	fmt.Println(string(body))
}

func TestCancleHttpConnect(t *testing.T)  {
	tr:=http.Transport{DisableKeepAlives: true}
	client:=http.Client{Transport: &tr}
	resp,err:=client.Get("https://golang.google.cn/")
	if resp!=nil{
		defer resp.Body.Close()
	}
	CheckError(err)
	fmt.Println(resp.StatusCode)
	body,err:=ioutil.ReadAll(resp.Body)
	CheckError(err)
	fmt.Println(len(body))
}