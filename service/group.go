package service

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prynnekey/gin-im/common"
	"github.com/prynnekey/gin-im/common/response"
	"github.com/prynnekey/gin-im/models"
)

// 创建群聊
func GroupCreateChat() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var roomBasicSimple models.RoomBasicSimple
		// 获取参数
		err := ctx.ShouldBind(&roomBasicSimple)
		if err != nil {
			ctx.JSON(http.StatusOK, response.Fail(nil, "参数格式不正确"))
			return
		}

		userIdentity := ctx.MustGet("user_claim").(*common.UserClaim).Identity

		// 生成房间号
		roomNum := common.GenerateCode(6)

		// 判断房间号是否存在
		_, err2 := models.GetRoomBasicByRoomNumber(roomNum)
		for i := 6; err2 == nil; i++ {
			roomNum = common.GenerateCode(i)
			_, err2 = models.GetRoomBasicByRoomNumber(roomNum)
		}

		// 初始化
		roomBasic := models.RoomBasic{
			Identity:     common.GenerateUUID(),
			UserIdentity: userIdentity,
			Number:       roomNum,
			Name:         roomBasicSimple.Name,
			Info:         roomBasicSimple.Info,
			CreateAt:     time.Now().Unix(),
			UpdateAt:     time.Now().Unix(),
		}

		// 插入RoomBasic
		if err := models.InsertOneRoomBasic(&roomBasic); err != nil {
			ctx.JSON(http.StatusOK, response.Fail(nil, "创建群聊时发生错误."+err.Error()))
			return
		}

		// 插入UserRoom
		userRoom := models.UserRoom{
			UserIdentity: userIdentity,
			RoomIdentity: roomBasic.Identity,
			RoomType:     2,
			CreateAt:     time.Now().Unix(),
			UpdateAt:     time.Now().Unix(),
		}

		if err := models.InsertOneUserRoom(&userRoom); err != nil {
			ctx.JSON(http.StatusOK, response.Fail(nil, "创建群聊时发生错误."+err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, response.Success(roomBasic, "创建成功"))
	}
}

func GroupGetCreateGroupChats() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		identity := ctx.MustGet("user_claim").(*common.UserClaim).Identity

		ur, err := models.GetUserRoomByUserIdentity(identity, 2)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusOK, response.Fail(nil, "获取群聊列表时发生错误."+err.Error()))
			return
		}

		rbs := make([]*models.RoomBasic, 0)
		for _, u := range ur {
			rb, err := models.GetRoomBasicByRoomIdentityAndUserIdentity(u.RoomIdentity, identity)
			if err != nil {
				log.Println(err)
				// ctx.JSON(http.StatusOK, response.Fail(nil, "获取群聊列表时发生错误."+err.Error()))
				continue
			}
			rbs = append(rbs, rb)
		}

		ctx.JSON(http.StatusOK, response.Success(rbs, "获取成功"))

	}
}

func GroupGetJoinedGroupChats() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		identity := ctx.MustGet("user_claim").(*common.UserClaim).Identity

		ur, err := models.GetUserRoomByUserIdentity(identity, 2)
		if err != nil {
			ctx.JSON(http.StatusOK, response.Fail(nil, "获取群聊列表时发生错误."+err.Error()))
			return
		}

		rbs := make([]*models.RoomBasic, 0)
		for _, u := range ur {
			rb, err := models.GetRoomBasicByRoomIdentity(u.RoomIdentity)
			if err != nil {
				ctx.JSON(http.StatusOK, response.Fail(nil, "获取群聊列表时发生错误."+err.Error()))
				return
			}
			rbs = append(rbs, rb)
		}

		ctx.JSON(http.StatusOK, response.Success(rbs, "获取成功"))
	}
}

func GroupInvateJoinedGroupChats() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取参数
		number := ctx.Param("number")

		// 校验参数
		if number == "" {
			ctx.JSON(http.StatusOK, response.Fail(nil, "请输入群号"))
			return
		}

		// 判断是否已经加入群聊
		rb, err := models.GetRoomBasicByRoomNumber(number)
		if err != nil {
			ctx.JSON(http.StatusOK, response.Fail(nil, "群号不存在"))
			return
		}

		identity := ctx.MustGet("user_claim").(*common.UserClaim).Identity
		_, err = models.GetUserRoomByUserIdentityAndRoomIdentityWithRoomType(identity, rb.Identity, 2)
		if err == nil {
			ctx.JSON(http.StatusOK, response.Fail(nil, "已加入群聊"))
			return
		}

		// 加入群聊
		userRoom := &models.UserRoom{
			UserIdentity: identity,
			RoomIdentity: rb.Identity,
			RoomType:     2,
			CreateAt:     time.Now().Unix(),
			UpdateAt:     time.Now().Unix(),
		}
		err2 := models.InsertOneUserRoom(userRoom)
		if err2 != nil {
			ctx.JSON(http.StatusOK, response.Fail(nil, "加入群聊失败"))
			return
		}

		ctx.JSON(http.StatusOK, response.Success(nil, "加入群聊成功"))
	}
}
