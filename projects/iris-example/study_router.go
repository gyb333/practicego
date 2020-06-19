package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

func api()  {
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())

	//GET请求: http://localhost:8000/api/getPath?name=davie&age=123 Body:none
	app.Get("/api/getPath", func(context context.Context) {
		path := context.Path()
		context.URLParams()
		app.Logger().Info(path,context.URLParams())
		name:=context.URLParam("name")
		pwd:=context.URLParam("age")
		data:=fmt.Sprintf("<h1>name:%s,age:%s</h1>", name, pwd)

		//context.WriteString(path)
		//context.Text(data)
		context.HTML(data)

	})

	/*2.处理Post请求 form表单的字段获取
	 url:POST http://localhost:8000/api/postForm?urlParam=test Body:Form-Data
			 name:gyb333
			 age:34
	 */
	app.Post("/api/postForm", func(context context.Context) {
		path := context.Path()
		app.Logger().Info(path)
		//context.PostValue方法来获取post请求所提交的for表单数据
		name := context.PostValue("name")
		pwd := context.PostValue("age")
		app.Logger().Info(name, "  ", pwd)
		context.HTML(name+" "+context.URLParam("urlParam"))
	})

	/**4、处理Post请求 Json格式数据
	 * Postman工具选择[{"key":"Content-Type","value":"application/json","description":""}]
	  请求内容：POST:http://localhost:8000/api/postJson?urlParam=test Body:raw
		{"name": "gyb333","age": 28}
	 */
	app.Post("/api/postJson", func(context context.Context) {

		//1.path
		path := context.Path()
		app.Logger().Info("请求URL：", path," URLParam: ",context.URLParam("urlParam"))

		//2.Json数据解析
		var person Person
		//context.ReadJSON()
		if err := context.ReadJSON(&person); err != nil {
			panic(err.Error())
		}

		//输出：Received: main.Person{Name:"davie", Age:28}
		context.Writef("Received: %#+v\n", person)
		context.JSON(person)

	})

	/**5.处理Post请求 Xml格式数据
	 * 请求配置：Content-Type到application/xml（可选但最好设置）
	  请求内容：POST http://localhost:8000/api/postXml?urlParam=test Body:raw
	   <student>
	 		<stu_name>davie</stu_name>
	 		<stu_age>28</stu_age>
	 	</student>
	 */
	app.Post("/api/postXml", func(context context.Context) {
		//1.Path
		path := context.Path()
		app.Logger().Info("请求URL：", path," URLParam: ",context.URLParam("urlParam"))

		//2.XML数据解析
		var student Student
		if err := context.ReadXML(&student); err != nil {
			panic(err.Error())
		}
		//输出：
		context.Writef("Received：%#+v\n", student)
		context.XML(student)
	})

	/*
		PUT: http://localhost:8000/api/put?urlParam=test	Body:Form-Data
	 		name:gyb333
			age: 34
	 */
	app.Put("/api/put", func(context context.Context) {
		path := context.Path()
		app.Logger().Info("请求url：", path," URLParam: ",context.URLParam("urlParam"))
		name := context.PostValue("name")
		pwd := context.PostValue("age")
		app.Logger().Info(name, "  ", pwd)
	})

	/*
		DELETE:http://localhost:8000/api/delete?urlParam=test Body:raw
		{"name": "gyb333","age": 28}
	 */
	app.Delete("/api/delete", func(context context.Context) {
		path := context.Path()
		app.Logger().Info("Delete请求url：", path," URLParam: ",context.URLParam("urlParam"))
		//2.Json数据解析
		var person Person
		//context.ReadJSON()
		if err := context.ReadJSON(&person); err != nil {
			panic(err.Error())
		}
		//输出：Received: main.Person{Name:"davie", Age:28}
		context.Writef("Received: %#+v\n", person)
	})


	//GET： http://localhost:8000/api/2019-03-10/beijing
	//      http://localhost:8000/api/2019-03-11/beijing
	//      http://localhost:8000/api/2019-03-11/tianjin
	app.Get("/api/{date}/{city}", func(context context.Context) {
		path := context.Path()
		date := context.Params().Get("date")
		city := context.Params().Get("city")
		context.WriteString(path + "  , " + date + " , " + city)
	})

	/* 1、Get 正则表达式 {name}路由
	   使用：context.Params().Get("name") 获取正则表达式变量
	   请求1：/hello/1  /hello/2  /hello/3 /hello/10000
	 */
	app.Get("/api/{name}", func(context context.Context) {
		//获取变量
		path := context.Path()

		app.Logger().Info(path)
		//获取正则表达式变量内容值
		name := context.Params().Get("name")
		context.HTML("<h1>" + name + "</h1>")
	})

	/* 2、自定义正则表达式变量路由请求 {unit64:uint64}进行变量类型限制
		GET: http://localhost:8000/api/uint/10
	 */
	app.Get("/api/uint/{userid:uint64}", func(context context.Context) {
		userID, err := context.Params().GetUint("userid")
		if err != nil {
			//设置请求状态码，状态码可以自定义
			context.JSON(map[string]interface{}{
				"requestcode": 201,
				"message":     "bad request",
			})
			return
		}
		context.JSON(map[string]interface{}{
			"requestcode": 200,
			"user_id":     userID,
		})
	})

	//自定义正则表达式路由请求 bool	http://localhost:8000/api/bool/false
	app.Get("/api/bool/{isLogin:bool}", func(context context.Context) {
		isLogin, err := context.Params().GetBool("isLogin")
		if err != nil {
			context.StatusCode(iris.StatusNonAuthoritativeInfo)
			return
		}
		if isLogin {
			context.WriteString(" 已登录 ")
		} else {
			context.WriteString(" 未登录 ")
		}
		//正则表达式所支持的数据类型
		context.Params()
	})



	app.Run(
		//开启web服务
		iris.Addr("localhost:8000"),
		// 禁用更新
		//iris.WithoutVersionChecker,
		// 按下CTRL / CMD + C时跳过错误的服务器：
		iris.WithoutServerError(iris.ErrServerClosed),
		//实现更快的json序列化和更多优化：
		iris.WithOptimizations,
	)
}

//自定义的struct
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

//自定义的结构体
type Student struct {
	//XMLName xml.Name `xml:"student"`
	Name string `xml:"name"`
	Age  int    `xml:"age"`
}
