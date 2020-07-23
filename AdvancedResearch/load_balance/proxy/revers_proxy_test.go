package proxy

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestReverse(t *testing.T)  {
	http.HandleFunc("/",handler)
	fmt.Println("Http reverse proxy server start at : 127.0.0.1:8888")
	if err := http.ListenAndServe(":8888",nil);err != nil{
		log.Fatal("Start server failed,err:",err)
	}
}
