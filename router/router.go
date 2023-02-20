package router

import (
	"github.com/gin-gonic/gin"
	"github.com/prynnekey/gin-im/service"
)

func Init() *gin.Engine {
	r := gin.Default()

	// 用户登录
	r.POST("/login", service.Login())
	// 用户注册
	r.POST("/register", service.Register())

	return r
}
