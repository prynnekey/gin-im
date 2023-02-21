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

		// 查询指定用户的个人信息
		auth.GET("/user/info/:username", service.UserInfo())

		// 添加好友
		auth.POST("/user/add/:username", service.AddUser())

		// 发送、接收消息
		auth.GET("/websocket/message", service.WebsocketMessage())

		// 获取聊天记录
		auth.GET("/chat/history", service.ChatHistory())
	}

	return r
}
