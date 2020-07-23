package httpServer

import (
	"log"
	"net/http"
	"testing"
)

func TestReveseProxy(t *testing.T)  {
	proxy := &ReveseProxyHandler{}

	log.Println("Start to serve at 127.0.0.1:8888")
	if err := http.ListenAndServe(":8888",proxy);err !=nil{
		log.Fatal("Failed to start reverse proxy server ,err:",err)
	}
}
