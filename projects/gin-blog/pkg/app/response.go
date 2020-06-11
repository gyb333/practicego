package app

import (
	"gin-blog/pkg/codec"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func  ResponseFunc(c *gin.Context,httpCode, errCode int, data interface{}) {
	c.JSON(httpCode,
		//gin.H{
		//"code": httpCode,
		//"msg":  codec.GetMsg(errCode),
		//"data": data,
		//},
		Response{
			Code: errCode,
			Msg:  codec.GetMsg(errCode),
			Data: data,
		},
	)

	return
}
