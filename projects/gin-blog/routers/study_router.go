package routers

import (
	"fmt"
	"gin-blog/pkg/logging"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)
type Person struct {
	Name string `form:"name"`
	Age  int    `form:"age"`
}
//自定义的结构体
type Student struct {
	Name string `xml:"name"`
	Age  int    `xml:"age"`
}


func apis(r *gin.Engine)  {
	api := r.Group("/api")
	{
		//GET请求: http://localhost:8000/api/getPath?name=davie&age=123 Body:none
		api.GET("/getPath", func(c *gin.Context) {
			path :=c.FullPath()
			logging.Info(path)
			var person Person
			// If `GET`, only `Form` binding engine (`query`) used.
			// If `POST`, first checks the `content-type` for `JSON` or `XML`, then uses `Form` (`form-data`).
			// See more at https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L48
			c.ShouldBind(&person)
			c.String(http.StatusOK,"Name:%s,Age:%d",person.Name,person.Age)
		})

		/*2.处理Post请求 form表单的字段获取
		 url:POST http://localhost:8000/api/postForm?urlParam=test Body:Form-Data
				 name:gyb333
				 age:11
		*/
		api.POST("/postForm", func(c *gin.Context) {
			fmt.Println(c.Query("urlParam"))
			var person Person
			c.ShouldBind(&person)
			c.String(http.StatusOK,"Name:%s,Age:%d",person.Name,person.Age)
		})

		/**4、处理Post请求 Json格式数据
		 * Postman工具选择[{"key":"Content-Type","value":"application/json","description":""}]
		  请求内容：POST:http://localhost:8000/api/postJson?urlParam=test Body:raw
			{"name": "gyb333","age": 28}
		*/
		api.POST("/postJson", func(c *gin.Context) {
			fmt.Println(c.Query("urlParam"))
			var person Person
			c.ShouldBind(&person)
			c.String(http.StatusOK,"Name:%s,Age:%d",person.Name,person.Age)
		})
		/**5.处理Post请求 Xml格式数据
		* 请求配置：Content-Type到application/xml（可选但最好设置）
		 请求内容：POST http://localhost:8000/api/postXml?urlParam=test Body:raw
		  	<student>
						<name>davie</name>
						<age>28</age>
			</student>
		*/
		api.POST("/postXml", func(c *gin.Context) {
			fmt.Println(c.Query("urlParam"))
			var student Student
			c.ShouldBind(&student)
			c.String(http.StatusOK,"Name:%s,Age:%d",student.Name,student.Age)
		})

		/*
				PUT: http://localhost:8000/api/put?urlParam=test	Body:Form-Data
			 		name:gyb333
					age: 23
		*/
		api.PUT("/put", func(c *gin.Context) {
			fmt.Println(c.Query("urlParam"))
			var person Person
			c.ShouldBind(&person)
			c.String(http.StatusOK,"Name:%s,Age:%d",person.Name,person.Age)
		})
		/*
			DELETE:http://localhost:8000/api/delete?urlParam=test Body:raw
			{"name": "gyb333","age": 28}
		*/
		api.DELETE("/delete", func(c *gin.Context) {
			fmt.Println(c.Query("urlParam"))
			var person Person
			c.ShouldBind(&person)
			c.String(http.StatusOK,"Name:%s,Age:%d",person.Name,person.Age)
		})


		//GET： http://localhost:8000/api/2019-03-10/beijing
		//      http://localhost:8000/api/2019-03-11/beijing
		//      http://localhost:8000/api/2019-03-11/tianjin

		/* 1、Get 正则表达式 {name}路由
		   使用：context.Params().Get("name") 获取正则表达式变量
		   请求1：http://localhost:8000/api/hello/110
		*/
		r.LoadHTMLGlob("view/*")
		api.GET("/hello/:name/*action", func(c *gin.Context) {
			//获取正则表达式变量内容值
			name := c.Param("name")
			action := c.Param("action")
			action = strings.Trim(action, "/")
			c.HTML(http.StatusOK,"index.html",gin.H{"title": "gin学习 "+action,"name": name})
		})
		/* 2、自定义正则表达式变量路由请求 {unit64:uint64}进行变量类型限制
		GET: http://localhost:8000/api/uint/10
		*/

		//自定义正则表达式路由请求 bool	http://localhost:8000/api/bool/false
	}
}
