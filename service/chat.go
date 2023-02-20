package service

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prynnekey/gin-im/common"
	"github.com/prynnekey/gin-im/common/response"
	"github.com/prynnekey/gin-im/models"
)

// 获取聊天记录
func ChatHistory() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取参数
		roomIdentity := ctx.Query("room_identity")
		pageIndex, _ := strconv.ParseInt(ctx.DefaultQuery("page_index", "1"), 10, 64)
		pageSize, err := strconv.ParseInt(ctx.DefaultQuery("page_size", "10"), 10, 64)

		// 参数校验
		if err != nil {
			ctx.JSON(http.StatusOK, response.Fail(nil, "参数类型错误"))
			return
		}

		if roomIdentity == "" {
			ctx.JSON(http.StatusOK, response.Fail(nil, "房间号不能为空"))
			return
		}

		// 获取user_claim
		userClaim := ctx.MustGet("user_claim").(*common.UserClaim)

		// 判断用户是否属于该房间
		_, err2 := models.GetUserRoomByUserIdentityAndRoomIdentity(userClaim.Identity, roomIdentity)
		if err2 != nil {
			log.Println("当前用户不属于该房间,无法查询聊天记录:", err2)
			ctx.JSON(http.StatusOK, response.Fail(nil, "当前用户不属于该房间,无法查询聊天记录"))
			return
		}

		// 根据房间identity查询聊天记录
		mbs, err := models.GetMessageBasicByRootIdentity(roomIdentity, pageIndex, pageSize)
		if err != nil {
			ctx.JSON(http.StatusOK, response.Fail(nil, "发生错误:"+err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, response.Success(gin.H{"history": mbs}, "查询成功"))
	}
}
