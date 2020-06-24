package routers

import (
	"gin-blog/middleware/cors"
	"gin-blog/middleware/jwt"
	"gin-blog/pkg/export"
	"gin-blog/pkg/qrcode"
	"gin-blog/pkg/upload"
	"gin-blog/routers/api"
	v1 "gin-blog/routers/api/v1"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"

	"gin-blog/pkg/setting"
	"github.com/gin-contrib/sessions"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())
	// 使用跨域中间件
	r.Use(cors.Cors())

	gin.SetMode(setting.RunMode)
	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	r.POST("/auth", api.GetAuth)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/upload", api.UploadImage)

	r.GET("/cookie", myCookie )

	//// 创建基于cookie的存储引擎，secret 参数是用于加密的密钥
	//store := cookie.NewStore([]byte("secret"))
	//// 设置session中间件，参数mysession，指的是session的名字，也是cookie的名字
	//r.Use(sessions.Sessions("mySession", store))

	store, _ := redis.NewStore(10, "tcp", "hadoop:6379", "", []byte("secret"))
	r.Use(sessions.Sessions("mySession", store))

	r.GET("/session", func(c *gin.Context) {
		// 初始化session对象
		session := sessions.Default(c)
		// 通过session.Get读取session值 session是键值对格式数据，因此需要通过key查询数据
		if session.Get("SessionID") != "gyb333" {
			// 设置session数据
			session.Set("SessionID", "gyb333")
			// 删除session数据
			session.Delete("tizi365")
			// 保存session数据
			session.Save()
			// 删除整个session
			// session.Clear()
		}

		c.JSON(200, gin.H{"hello": session.Get("SessionID")})
	})

	r.GET("/incr", func(c *gin.Context) {
		session := sessions.Default(c)
		var count int
		v := session.Get("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count++
		}
		session.Set("count", count)
		session.Save()
		c.JSON(200, gin.H{"count": count})
	})
	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{
		//获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		//新建标签
		apiv1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)

		//导出标签
		r.POST("/tags/export", v1.ExportTag)
		//导入标签
		r.POST("/tags/import", v1.ImportTag)

		//获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		//获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		//新建文章
		apiv1.POST("/articles", v1.AddArticle)
		//更新指定文章
		apiv1.PUT("/articles/:id", v1.EditArticle)
		//删除指定文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
		//生成文章海报
		apiv1.POST("/articles/poster/generate", v1.GenerateArticlePoster)
	}


	apis(r)

	return r
}

// @Summary 测试
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/cookie [get]
func myCookie(c *gin.Context) {
	// 根据cookie名字读取cookie值
	data, err := c.Cookie("my_cookie")
	if err != nil {
		data="cookie_value"
		//Secure=true，那么这个cookie只能用https协议发送给服务器
		//通过将cookie的MaxAge设置为-1, 达到删除cookie的目的。
		c.SetCookie("my_cookie",data,3600,"/","localhost",false,true)
	}
	c.JSON(200, gin.H{"code":200,"data":"Hello world!","msg":data+" ok"})
}