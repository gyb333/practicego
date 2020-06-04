package initRouter

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/testdata/protoexample"
	"net/http"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	// 添加 Get 请求路由
	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "hello gin")
	})

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "views/404.html", nil)
	})
	//	/user/ls ，和 /user/hello 都可以匹配，⽽ /user/ 和 /user/ls/ 不会被匹配
	router.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

	router.GET("/welcome", func(c *gin.Context) {
		name := c.DefaultQuery("name", "Guest") //可设置默认值
		//nickname := c.Query("nickname") // 是
		c.Request.URL.Query().Get("nickname")
		c.String(http.StatusOK, fmt.Sprintf("Hello %s ", name))
	})

	/*
	http的报文体传输数据就比query string稍微复杂⼀点，常见的格式就有四种。例如
			application/json ， application/x-www-form-urlencoded , application/xml 和multipart/form-data 。
		multipart/form-data：主要用于图片上传。json格式的很好理解，urlencode其实也不难，⽆⾮就是把query string的内容，放到了body体⾥，同样也需要urlencode。默认情况下，
		c.PostFROM解析的是 x-www-form-urlencoded 或 from-data 的参数
	 */
	router.POST("/form", func(c *gin.Context) {
		type1 := c.DefaultPostForm("type", "alert") //可设置默认值
		username := c.PostForm("username")
		password := c.PostForm("password")
		//hobbys := c.PostFormMap("hobby")
		//hobbys := c.QueryArray("hobby")
		hobbys := c.PostFormArray("hobby")
		c.String(http.StatusOK, fmt.Sprintf("type is %s, username is %s,password is %s,hobby is %v",
			type1, username, password,hobbys))
	})

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	//router.Static("/", "./public")
	router.POST("/upload", func(c *gin.Context) {
		// Multipart form
		form, err := c.MultipartForm()
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s",
				err.Error()))
			return
		}
		files := form.File["files"]
		for _, file := range files {
			if err := c.SaveUploadedFile(file, file.Filename); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
				return
			}
		}
		c.String(http.StatusOK, fmt.Sprintf("Uploaded successfully %d files ",
			len(files)))
	})


	// Simple group: v1
	v1 := router.Group("/v1")
	{
		v1.GET("/login", loginEndpoint)
		v1.GET("/submit", submitEndpoint)
		v1.POST("/read", readEndpoint)
	}
	// Simple group: v2
	v2 := router.Group("/v2")
	{
		v2.POST("/login", loginEndpoint)
		v2.POST("/submit", submitEndpoint)
		v2.POST("/read", readEndpoint)
	}

	router.POST("/loginJSON", func(c *gin.Context) {
		var json Login
		//其实就是将request中的Body中的数据按照JSON格式解析到json变量中
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if json.UserName != "ls" || json.Password != "123456" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})


	// Example for binding a HTML form (user=ls&password=123456)
	router.POST("/loginForm", func(c *gin.Context) {
		var form Login
		//⽅法⼀：对于FORM数据直接使⽤Bind函数, 默认使⽤使⽤form格式解析,if c.Bind(&form) == nil
		// 根据请求头中 content-type ⾃动推断.
		if err := c.Bind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if form.UserName != "ls" || form.Password != "123456" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})
	router.POST("/login", func(c *gin.Context) {
		var form Login
		//⽅法⼆: 使⽤BindWith函数,如果你明确知道数据的类型
		// 你可以显式声明来绑定多媒体表单：
		// c.BindWith(&form, binding.Form)
		// 或者使⽤⾃动推断:
		if c.BindWith(&form, binding.Form) == nil {
			if form.UserName == "user" && form.Password == "password" {
				c.JSON(200, gin.H{"status": "you are logged in ..... "})
			} else {
				c.JSON(401, gin.H{"status": "unauthorized"})
			}
		}
	})

	// URI
	//router.GET("/:username/:password", func(c *gin.Context) {
	//	var login Login
	//	if err := c.ShouldBindUri(&login); err != nil {
	//		c.JSON(400, gin.H{"msg": err})
	//		return
	//	}
	//	c.JSON(200, gin.H{"username": login.UserName, "password": login.Password})
	//})


	//既然请求可以使⽤不同的 content-type ，响应也如此。通常响应会有html，text，plain，json和xml等。 Gin提供了很优雅的渲染⽅法

	// gin.H is a shortcut for map[string]interface{}
	router.GET("/someJSON", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})
	router.GET("/moreJSON", func(c *gin.Context) {
		// You also can use a struct
		var msg struct {
			Name string `json:"user"`
			Message string
			Number int
		}
		msg.Name = "ls"
		msg.Message = "hey"
		msg.Number = 123
		// 注意 msg.Name 变成了 "user" 字段
		// 以下⽅式都会输出 : {"user": "ls", "Message": "hey", "Number": 123}
		c.JSON(http.StatusOK, msg)
	})
	router.GET("/someXML", func(c *gin.Context) {
		c.XML(http.StatusOK, gin.H{"user":"ls","message": "hey", "status":
		http.StatusOK})
	})
	router.GET("/someYAML", func(c *gin.Context) {
		c.YAML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})
	router.GET("/someProtoBuf", func(c *gin.Context) {
		reps := []int64{int64(1), int64(2)}
		label := "test"
		// The specific definition of protobuf is written in the testdata/protoexample file.
			data := &protoexample.Test{
			Label: &label,
			Reps: reps,
		}
		// Note that data becomes binary data in the response
		// Will output protoexample.Test protobuf serialized data
		c.ProtoBuf(http.StatusOK, data)
	})

	return router
}


func loginEndpoint(c *gin.Context) {
	name := c.DefaultQuery("name", "Guest") //可设置默认值
	c.String(http.StatusOK, fmt.Sprintf("Hello %s \n", name))
}
func submitEndpoint(c *gin.Context) {
	name := c.DefaultQuery("name", "Guest") //可设置默认值
	c.String(http.StatusOK, fmt.Sprintf("Hello %s \n", name))
}
func readEndpoint(c *gin.Context) {
	name := c.DefaultQuery("name", "Guest") //可设置默认值
	c.String(http.StatusOK, fmt.Sprintf("Hello %s \n", name))
}

type Login struct {
	UserName string `form:"username" json:"username" uri:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" uri:"password" xml:"password" binding:"required"`
}
