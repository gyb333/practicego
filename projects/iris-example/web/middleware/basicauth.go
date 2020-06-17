package middleware

import "github.com/kataras/iris/middleware/basicauth"

// 简单的授权验证
var BasicAuth = basicauth.New(basicauth.Config{
	Users: map[string]string{
		"admin": "password",
	},
})
