package group

import (
	"github.com/gin-gonic/gin"
	"github.com/prynnekey/gin-im/service"
)

func InitRouter(auth *gin.RouterGroup) {
	group := auth.Group("/group")
	{
		// 创建群聊
		group.POST("/create_chat", service.GroupCreateChat())

		// 查看我创建的群聊
		group.GET("/create_chat", service.GroupGetCreateGroupChats())

		// 查看已加入的群聊
		group.GET("/joined", service.GroupGetJoinedGroupChats())

		// 加入群聊
		group.POST("/invite/:number", service.GroupInvateJoinedGroupChats())

		// 移除用户从群聊

		// 退出指定群聊

		// 解散群聊

		// 将群管理员移交给他人

		// 查看群中所有用户
		group.GET("/all_user/:number", service.GroupGetAllUser())

		// 获取群聊详细信息
	}

}
