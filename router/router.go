package router

import (
	"github.com/gin-gonic/gin"
	"github.com/prynnekey/gin-im/middleware"
	"github.com/prynnekey/gin-im/service"
)

func Init() *gin.Engine {
	r := gin.Default()

	// 用户登录
	r.POST("/login", service.Login())
	// 用户注册
	r.POST("/register", service.Register())

	// 需要登录才能访问的接口
	auth := r.Group("/api", middleware.AuthHandler())
	{
		// 获取用户详情
		auth.GET("/user/detail", service.UserDetail())

		// 发送、接收消息
		auth.GET("/websocket/message", service.WebsocketMessage())
	}

	return r
}
