package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/mvc"
	"iris-example/web/controllers"
	"sync/atomic"
)

func main() {
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())

	app.Logger().SetLevel("debug")
	//加载模板文件
	app.RegisterView(iris.HTML("./web/views", ".html").
		Layout("shared/layout.html").
		Reload(true))
	app.StaticWeb("/static", "./web/static")

	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().
			GetStringDefault("message", "The page you're looking for doesn't exist"))
		ctx.View("shared/error.html")
	})

	// 注册控制器
	mvc.New(app.Party("/hello")).Handle(new(controllers.HelloController))
	mvc.New(app.Party("/visits")).Handle(&globalVisitorsController{visits: 0})
	mvc.New(app).Handle(new(controllers.ExampleController))
	//您还可以拆分您编写的代码以配置mvc.Application
	//使用`mvc.Configure`方法，如下所示。
	mvc.Configure(app.Party("/movies"), movies)
	// http://localhost:8080/movies
	// http://localhost:8080/movies/1

	mvc.Configure(app.Party("/users"),users)
	mvc.Configure(app.Party("/user"),user)

	mvc.Configure(app.Party("/websocket"), configureMVC)

	app.Run(
		//开启web服务
		iris.Addr("localhost:8080"),
		// 禁用更新
		//iris.WithoutVersionChecker,
		// 按下CTRL / CMD + C时跳过错误的服务器：
		iris.WithoutServerError(iris.ErrServerClosed),
		//实现更快的json序列化和更多优化：
		iris.WithOptimizations,
	)
}

type globalVisitorsController struct {
	//当使用单例控制器时，由开发人员负责访问安全,所有客户端共享相同的控制器实例。
	//注意任何控制器的方法,是每个客户端，但结构的字段可以在多个客户端共享（如果是结构）
	//没有任何依赖于的动态struct字段依赖项
	//并且所有字段的值都不为零，在这种情况下我们使用uint64，它不是零（即使我们没有设置它手动易于理解的原因）因为它的值为＆{0}
	//以上所有都声明了一个Singleton，请注意，您不必编写一行代码来执行此操作，Iris足够聪明。
	//见`Get`
	visits uint64
}

func (c *globalVisitorsController) Get() string {
	count := atomic.AddUint64(&c.visits, 1)
	return fmt.Sprintf("Total visitors: %d", count)
}