package user

import (
	"github.com/gin-gonic/gin"
	"github.com/prynnekey/gin-im/service"
)

func InitRouter(auth *gin.RouterGroup) {
	user := auth.Group("/user")
	{
		// 获取用户详情
		user.GET("/detail", service.UserDetail())

		// 查询指定用户的个人信息
		user.GET("/info/:username", service.UserInfo())

		// 添加好友
		user.POST("/add/:username", service.UserAdd())

		// 删除好友
		user.DELETE("/delete/:username", service.UserDelete())

		// 查询我的所有好友
	}
}
