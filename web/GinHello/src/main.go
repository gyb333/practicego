package main

import (
	_ "GinHello/src/docs"
	"GinHello/src/initRouter"
	"fmt"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
	"time"
)
//全局中间件
func MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		fmt.Println("before middleware")
		//设置request变量到Context的Key中,通过Get等函数可以取得
		c.Set("request", "client_request")
		//发送request之前
		c.Next()
		//发送requst之后
		// 这个c.Write是ResponseWriter,我们可以获得状态等信息
		status := c.Writer.Status()
		fmt.Println("after middleware,", status)
		t2 := time.Since(t)
		fmt.Println("time:", t2)
	}
}


// @title Golang Esign API
// @version 1.0
// @description  Golang api of ginweb
// @termsOfService http://github.com
// @contact.name API Support
// @contact.url http://www.cnblogs.com
// @contact.email 876368840@qq.com
//@host localhost:8080
func main() {
	router := initRouter.SetupRouter()
	// 显示当前⽂件夹下的所有⽂件/或者指定⽂件
	router.StaticFS("/showDir", http.Dir("."))
	router.StaticFS("/files", http.Dir("/bin"))
	//Static提供给定⽂件系统根⽬录中的⽂件。
	//router.Static("/files", "/bin")
	router.StaticFile("/image", "./assets/image.jpg")
	//http重定向
	router.GET("/redirect", func(c *gin.Context) {
		//⽀持内部和外部的重定向
		c.Redirect(http.StatusMovedPermanently, "http://www.baidu.com/")
	})
	//路由重定向
	router.GET("/test", func(c *gin.Context) {
		// 指定重定向的URL
		c.Request.URL.Path = "/test2"
		router.HandleContext(c)
	})
	router.GET("/test2", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"hello": "world"})
	})

	router.Use(MiddleWare())
	{
		router.GET("/middleware", func(c *gin.Context) {
			//获取gin上下⽂中的变量
			request := c.MustGet("request").(string)
			req, _ := c.Get("request")
			fmt.Println("request:",request)
			c.JSON(http.StatusOK, gin.H{
				"middile_request": request,
				"request": req,
			})
		})
	}

	//单个路由中间件
	router.GET("/before", MiddleWare(), func(c *gin.Context) {
		request := c.MustGet("request").(string)
		c.JSON(http.StatusOK, gin.H{
			"middile_request": request,
		})
	})

	//后使⽤ gin.BasicAuth 中间件，设置授权⽤户
	authorized := router.Group("/admin", gin.BasicAuth(gin.Accounts{
		"ls": "123",
		"yang": "111",
		"edu": "666",
		"lucy": "4321",
	}))
	//定义路由
	authorized.GET("/secrets", func(c *gin.Context) {
		// 获取提交的⽤户名（AuthUserKey）
		user := c.MustGet(gin.AuthUserKey).(string)
		if secret, ok := secrets[user]; ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET : ("})
			}
		})


	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	_ = router.Run()

}
// 模拟私有数据
var secrets = gin.H{
	"ls": gin.H{"email": "ls@lianshiclass.com", "phone": "123456"},
	"yang": gin.H{"email": "yang@lianshiclass.com", "phone": "111111"},
	"edu": gin.H{"email": "edu@lianshiclass.com", "phone": "666666"},
}


