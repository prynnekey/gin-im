package service

import (
	"flag"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/prynnekey/gin-im/common"
	"github.com/prynnekey/gin-im/define"
	"github.com/prynnekey/gin-im/models"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

var clients = make(map[string]*websocket.Conn) // connected clients

func WebsocketMessage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Print("upgrade:", err)
			// ctx.JSON(http.StatusOK, response.Fail(nil, err.Error()))
			return
		}
		defer c.Close()

		// 获取user_claim
		userClaim := ctx.MustGet("user_claim").(*common.UserClaim)

		// 将identity和连接放入map中
		clients[userClaim.Identity] = c

		for {
			ms := &define.Message{}
			// 读取客户端发送的数据
			err := c.ReadJSON(ms)
			if err != nil {
				log.Println("Read error:", err)
				// ctx.JSON(http.StatusOK, response.Fail(nil, err.Error()))
				return
			}

			// 判断用户是否属于这个房间
			_, err2 := models.GetUserRoomByUserIdentityAndRoomIdentity(userClaim.Identity, ms.RoomIdentity)
			if err2 != nil {
				log.Println("当前用户不属于该房间:", err2)
				// ctx.JSON(http.StatusOK, response.Fail(nil, "当前用户不属于该房间"))
				return
			}

			// 保存消息
			mb := models.MessageBasic{
				Identity:     common.GenerateUUID(),
				UserIdentity: userClaim.Identity,
				RootIdentity: ms.RoomIdentity,
				Data:         ms.Message,
				CreateAt:     time.Now().Unix(),
				UpdateAt:     time.Now().Unix(),
			}
			err = models.InsertOneMessageBasic(&mb)
			if err != nil {
				log.Println("保存消息失败:", err)
			}

			// 获取房间内的所有用户，然后发送消息
			userRooms, err := models.GetUserRoomByRoomIdentity(ms.RoomIdentity)
			if err != nil {
				log.Println("获取房间内用户失败:", err)
				// ctx.JSON(http.StatusOK, response.Fail(nil, "获取房间内用户失败"))
				return
			}

			// 发送消息
			for _, room := range userRooms {
				if c, ok := clients[room.UserIdentity]; ok {
					err := c.WriteMessage(websocket.TextMessage, []byte(ms.Message))
					if err != nil {
						log.Println("发送消息失败:", err)
						// ctx.JSON(http.StatusOK, response.Fail(nil, err.Error()))
						return
					}
				}
			}
		}

	}
}
