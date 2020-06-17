package main

import (
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"iris-example/datasource"
	"iris-example/repositories"
	"iris-example/services"
	"iris-example/web/controllers"
	"iris-example/web/middleware"
	"log"
	"time"
)



//注意mvc.Application，它不是iris.Application。
func movies(app *mvc.Application) {
	//添加基本身份验证（admin：password）中间件
	app.Router.Use(middleware.BasicAuth)
	// 使用数据源中的一些（内存）数据创建我们的电影资源库。
	repo := repositories.NewMovieRepository(datasource.Movies)
	// 创建我们的电影服务，我们将它绑定到电影应用程序的依赖项。
	movieService := services.NewMovieService(repo)
	app.Register(movieService)
	//为我们的电影控制器服务 ，您可以为多个控制器提供服务
	//你也可以使用`movies.Party（relativePath）`或`movies.Clone（app.Party（...））创建子mvc应用程序
	app.Handle(new(controllers.MovieController))
}

func init()  {
	// Prepare our repositories and services.
	db, err := datasource.LoadUsers(datasource.Memory)
	if err != nil {
		log.Fatalf("error while loading the users: %v", err)
		return
	}
	repo := repositories.NewUserRepository(db)
	userService = services.NewUserService(repo)
}
var userService services.UserService

func users(app *mvc.Application) {
	app.Router.Use(middleware.BasicAuth)
	app.Register(userService)
	app.Handle(new(controllers.UsersController))
}


func user(app *mvc.Application) {
	// "/user" based mvc application.
	sessManager := sessions.New(sessions.Config{
		Cookie:  "sessioncookiename",
		Expires: 24 * time.Hour,
	})

	app.Register(
		userService,
		sessManager.Start,
	)
	app.Handle(new(controllers.UserController))
}

