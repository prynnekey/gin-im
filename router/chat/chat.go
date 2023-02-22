package chat

import (
	"github.com/gin-gonic/gin"
	"github.com/prynnekey/gin-im/service"
)

func InitRouter(auth *gin.RouterGroup) {
	// 发送、接收消息
	auth.GET("/websocket/message", service.WebsocketMessage())

	// 获取聊天记录
	auth.GET("/chat/history", service.ChatHistory())
}
