package router

import (
	"github.com/gin-gonic/gin"
	"github.com/prynnekey/gin-im/middleware"
	"github.com/prynnekey/gin-im/router/chat"
	"github.com/prynnekey/gin-im/router/group"
	"github.com/prynnekey/gin-im/router/user"
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
		// 用户路由
		user.InitRouter(auth)

		// 群聊路由
		group.InitRouter(auth)

		// 聊天路由
		chat.InitRouter(auth)
	}

	return r
}
