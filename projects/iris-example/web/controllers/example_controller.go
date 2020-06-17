package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

// ExampleController提供 ”/”，“/ping”和 “/hello”路由选项
type ExampleController struct{}

// Get 服务
// 请求方法:   GET
// 请求资源路径: http://localhost:8080
func (c *ExampleController) Get() mvc.Result {
	return mvc.Response{
		ContentType: "text/html",
		Text:        "<h1>Welcome</h1>",
	}
}

// GetPing 服务
// 请求方法:   GET
// 请求资源路径: http://localhost:8080/ping
func (c *ExampleController) GetPing() string {
	return "pong"
}

// GetHello 服务
// 请求方法:   GET
// 请求资源路径: http://localhost:8080/hello
func (c *ExampleController) GetHi() interface{} {
	return map[string]string{"message": "Hello Iris!"}
}


func (c *ExampleController) BeforeActivation(b mvc.BeforeActivation) {
	anyMiddlewareHere := func(ctx iris.Context) {
		ctx.Application().Logger().Warnf("BeforeActivation")
		ctx.Next()
	}
	b.Handle("GET", "/custom_path",
		"CustomHandlerWithoutFollowingTheNamingGuide", anyMiddlewareHere)
	//甚至添加基于此控制器路由的全局中间件，
	//在这个例子中是根“/”：
	//b.Router().Use(middleware)
}

func (c *ExampleController) AfterActivation(a mvc.AfterActivation) {
	//if a.Singleton() {
	//	panic("basicController should be stateless,a request-scoped,we have a 'Session' which depends on the context.")
	//}

	//根据您想要的方法名称选择路由 修改
	index := a.GetRoute("CustomHandlerWithoutFollowingTheNamingGuide")
	//只是将处理程序作为您想要使用的中间件预先添加。
	anyMiddlewareHere := func(ctx iris.Context) {
		ctx.Application().Logger().Warnf("AfterActivation")
		ctx.Next()
	}
    index.Handlers = append([]iris.Handler{anyMiddlewareHere}, index.Handlers...)

}
// CustomHandlerWithoutFollowingTheNamingGuide 服务
// 请求方法:   GET
// 请求资源路径: http://localhost:8080/custom_path
func (c *ExampleController) CustomHandlerWithoutFollowingTheNamingGuide() string {
	return "hello from the custom handler without following the naming guide"
}

// GetUserBy 服务
// 请求方法:   GET
// 请求资源路径: http://localhost:8080/user/{username:string}
//是一个保留的关键字来告诉框架你要在函数的输入参数中绑定路径参数，
//在同一控制器中使用“Get”和“GetBy”可以实现
//
//func (c *ExampleController) GetUserBy(username string) mvc.Result {
//    return mvc.View{
//        Name: "user/username.html",
//        Data: username,
//    }
//}

/*
func (c *ExampleController) Post() {}
func (c *ExampleController) Put() {}
func (c *ExampleController) Delete() {}
func (c *ExampleController) Connect() {}
func (c *ExampleController) Head() {}
func (c *ExampleController) Patch() {}
func (c *ExampleController) Options() {}
func (c *ExampleController) Trace() {}
*/

/*
func (c *ExampleController) All() {}
//  或者
func (c *ExampleController) Any() {}

func (c *ExampleController) BeforeActivation(b mvc.BeforeActivation) {
    // 1 -> http 请求方法
    // 2 -> 请求路由
    // 3 -> 此控制器的方法名称应该是该路由的处理程序
    b.Handle("GET", "/mypath/{param}", "DoIt", optionalMiddlewareHere...)
}

//AfterActivation，所有依赖项都被设置,因此访问它们是只读
，但仍可以添加自定义控制器或简单的标准处理程序。
func (c *ExampleController) AfterActivation(a mvc.AfterActivation) {}
*/
