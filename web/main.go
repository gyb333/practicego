package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Web!")
}


type GreetingHandler struct {
	Language string
}

func (h GreetingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", h.Language)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	//x-www-form-urlencoded 编码，这是表单的默认编码。
	fmt.Fprint(w, `
<html>
 <head>
 <title>Go Web</title>
 </head>
 <body>
 <form method="post" action="/body">
 <label for="username">⽤户名：</label>
 <input type="text" id="username" name="username">
 <label for="email">邮箱：</label>
 <input type="text" id="email" name="email">
 <button type="submit">提交</button>
 </form>
 </body>
</html>
`)
}
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<html>
 <head>
 <title>Go Web</title>
 </head>
 <body>
 <form action="/form?lang=cpp&name=ls" method="post" enctype="application/x-www-form-urlencoded">
 <label>Form:</label>
 <input type="text" name="lang" />
 <input type="text" name="age" />
 <button type="submit">提交</button>
 </form>
 </body>
</html>`)
}
func worldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<html>
 <head>
 <title>Go Web</title>
 </head>
 <body>
 <form action="/multipartform?lang=cpp&name=dj" method="post" enctype="multipart/form-data">
 <label>MultipartForm:</label>
 <input type="text" name="lang" />
 <input type="text" name="age" />
 <input type="file" name="uploaded" />
 <button type="submit">提交</button>
</form>
 </body>
</html>`)
}
func urlHandler(w http.ResponseWriter, r *http.Request) {
	URL := r.URL

	fmt.Fprintf(w, "Scheme: %s\n", URL.Scheme)
	fmt.Fprintf(w, "Host: %s\n", URL.Host)
	fmt.Fprintf(w, "Path: %s\n", URL.Path)
	fmt.Fprintf(w, "RawPath: %s\n", URL.RawPath)
	fmt.Fprintf(w, "RawQuery: %s\n", URL.RawQuery)
	fmt.Fprintf(w, "Fragment: %s\n", URL.Fragment)
}
func protoFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Proto: %s\n", r.Proto)
	fmt.Fprintf(w, "ProtoMajor: %d\n", r.ProtoMajor)
	fmt.Fprintf(w, "ProtoMinor: %d\n", r.ProtoMinor)
}
func headerHandler(w http.ResponseWriter, r *http.Request) {
	for key, value := range r.Header {
		fmt.Fprintf(w, "%s: %v\n", key, value)
	}
}

func bodyHandler(w http.ResponseWriter, r *http.Request) {
	data := make([]byte, r.ContentLength)
	r.Body.Read(data) // 忽略错误处理
	defer r.Body.Close()

	//data, _ := ioutil.ReadAll(r.Body)

	fmt.Fprintln(w, string(data))
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	//以 localhost:8080/query?name=ls&age=20 请求，查询字符串 name=ls&age=20
	fmt.Fprintln(w, r.URL.RawQuery)
}

//enctype 指定请求体的编码⽅式，默认为 application/x-www-form-urlencoded 。如果需要发送⽂件，必须指定为 multipart/form-data ；
func formHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Fprintln(w, r.Form)
}

func multipartFormHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024)
	fmt.Fprintln(w, r.MultipartForm)

	fileHeader := r.MultipartForm.File["uploaded"][0]
	file, err := fileHeader.Open()
	if err != nil {
		fmt.Println("Open failed: ", err)
		return
	}
	data, err := ioutil.ReadAll(file)
	if err == nil {
		fmt.Fprintln(w, string(data))
	}
}
func writeHandler(w http.ResponseWriter, r *http.Request) {
	str := `<html>
<head><title>Go Web</title></head>
<body><h1>直接使⽤ Write ⽅法<h1></body>
</html>`
	w.Write([]byte(str))
}

func writeHeaderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "http://baidu.com")
	w.WriteHeader(501)
	fmt.Fprintln(w, "This API not implemented!!!")
}

func main() {
	//http.HandleFunc("/", hello)
	//if err := http.ListenAndServe(":8080", nil); err != nil {
	//	log.Fatal(err)
	//}
	//创建Mux
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)

	mux.Handle("/chinese", GreetingHandler{Language: "你好"})
	mux.Handle("/english", GreetingHandler{Language: "Hello"})

	mux.HandleFunc("/index", indexHandler)
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/hello/world", worldHandler)

	mux.HandleFunc("/url", urlHandler)

	mux.HandleFunc("/proto", protoFunc)

	mux.HandleFunc("/header", headerHandler)

	mux.HandleFunc("/body", bodyHandler)

	mux.HandleFunc("/query", queryHandler)

	mux.HandleFunc("/form", formHandler)

	mux.HandleFunc("/multipartform", multipartFormHandler)

	mux.HandleFunc("/write", writeHandler)

	mux.HandleFunc("/writeheader", writeHeaderHandler)

	server := &http.Server{
		Addr: ":8080",
		Handler: mux, //注册处理器
		ReadTimeout: 1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
