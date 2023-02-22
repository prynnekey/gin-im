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

// 用户登录
func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取参数
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")

		// 验证参数
		if username == "" || password == "" {
			ctx.JSON(http.StatusOK, response.Fail(nil, "参数错误"))
			return
		}

		// 将密码进行md5加密
		password = common.MD5(password)

		// 查询数据
		ub, err := models.GetUserBasicByUsernameAndPassword(username, password)
		if err != nil {
			ctx.JSON(http.StatusOK, response.Fail(nil, "用户名或密码错误"))
			return
		}

		// 登录成功
		// 生成token
		token, err := common.GenerateToken(ub.Identity, ub.Username)
		if err != nil {
			ctx.JSON(http.StatusOK, response.Fail(nil, "生成token时失败"+err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, response.Success(gin.H{"token": token}, "登录成功"))
	}
}

// 用户注册
func Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取参数
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")
		confirmPassword := ctx.PostForm("confirm_password")

		// 验证参数
		if username == "" {
			ctx.JSON(http.StatusOK, response.Fail(nil, "用户名不能为空"))
			return
		}

		if password == "" {
			ctx.JSON(http.StatusOK, response.Fail(nil, "密码不能为空"))
			return
		}

		if confirmPassword == "" {
			ctx.JSON(http.StatusOK, response.Fail(nil, "确认密码不能为空"))
			return
		}

		if password != confirmPassword {
			ctx.JSON(http.StatusOK, response.Fail(nil, "两次密码不一致"))
			return
		}

		// 判断用户名是否存在
		count, err := models.GetUserBasicCountByUsername(username)
		if err != nil {
			log.Printf("查询用户失败: %v", err)
			ctx.JSON(http.StatusOK, response.Fail(nil, "发生错误"+err.Error()))
			return
		}

		if count > 0 {
			ctx.JSON(http.StatusOK, response.Fail(nil, "用户名已存在"))
			return
		}

		// 将密码进行md5加密
		password = common.MD5(password)

		// 保存数据
		ub := models.UserBasic{
			Identity: common.GenerateUUID(),
			Username: username,
			Password: password,
			CreateAt: time.Now().Unix(),
			UpdateAt: time.Now().Unix(),
		}
		err = models.InsertOneUserBasic(&ub)
		if err != nil {
			log.Printf("插入数据失败: %v", err)
			ctx.JSON(http.StatusOK, response.Fail(nil, "发生错误"+err.Error()))
			return
		}

		// 注册成功

		// 生成token
		token, err := common.GenerateToken(ub.Identity, ub.Username)
		if err != nil {
			ctx.JSON(http.StatusOK, response.Fail(nil, "生成token时失败"+err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, response.Success(gin.H{"token": token}, "注册成功"))
	}
}

// 获取用户详情
func UserDetail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取用户claim
		value, exists := ctx.Get("user_claim")
		if !exists {
			ctx.JSON(http.StatusOK, response.Fail(nil, "获取用户信息失败"))
			return
		}

		// 解析claim,通过类型断言
		uc := value.(*common.UserClaim)

		// 查询数据
		ub, err := models.GetUserBasicByIdentity(uc.Identity)
		if err != nil {
			log.Printf("查询用户失败: %v", err)
			ctx.JSON(http.StatusOK, response.Fail(nil, "发生错误"+err.Error()))
			return
		}

		// 返回
		ctx.JSON(http.StatusOK, response.Success(gin.H{"user": ub}, "获取用户信息成功"))
	}
}

func UserInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取参数
		username := ctx.Param("username")

		// 参数校验
		if username == "" {
			ctx.JSON(http.StatusOK, response.Fail(nil, "要查询的用户名能为空"))
			return
		}

		// 查询数据
		ui, err := models.GetUserInfoByUsername(username)
		if err != nil {
			ctx.JSON(http.StatusOK, response.Fail(nil, "用户不存在"))
			return
		}

		userClaim, exists := ctx.Get("user_claim")
		if !exists {
			ctx.JSON(http.StatusOK, response.Fail(nil, "获取用户信息失败"))
			return
		}

		// 判断是否是好友关系
		isFirend, err := models.JudgeUserIsFriend(userClaim.(*common.UserClaim).Identity, ui.Identity)
		if err != nil {
			log.Printf("查询用户失败: %v", err)
			ctx.JSON(http.StatusOK, response.Fail(nil, "获取用户信息失败"))
			return
		}

		ui.IsFriend = isFirend

		ctx.JSON(http.StatusOK, response.Success(gin.H{"user": ui}, "获取用户信息成功"))
	}
}

// 添加好友
//
// 无需对方同意即可添加好友
func UserAdd() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取参数
		username := ctx.Param("username")

		// 参数校验
		if username == "" {
			ctx.JSON(http.StatusOK, response.Fail(nil, "用户名不能为空"))
			return
		}

		// 判断是否为好友
		ub, err := models.GetUserBasicByUsername(username)
		if err != nil {
			log.Printf("数据查询一场 err: %v\n", err)
			ctx.JSON(http.StatusOK, response.Fail(nil, "数据查询异常"))
			return
		}

		userClaim, exists := ctx.Get("user_claim")
		if !exists {
			ctx.JSON(http.StatusOK, response.Fail(nil, "获取用户信息失败"))
			return
		}

		isFirend, err := models.JudgeUserIsFriend(userClaim.(*common.UserClaim).Identity, ub.Identity)
		if err != nil {
			log.Printf("查询用户失败: %v", err)
			ctx.JSON(http.StatusOK, response.Fail(nil, "获取用户信息失败"))
			return
		}

		if isFirend {
			ctx.JSON(http.StatusOK, response.Fail(nil, "已经是好友了，无需添加"))
			return
		}

		// 添加好友

		// 1. 插入RoomBasic
		rb := models.RoomBasic{
			Identity:     common.GenerateUUID(),
			UserIdentity: userClaim.(*common.UserClaim).Identity,
			CreateAt:     time.Now().Unix(),
			UpdateAt:     time.Now().Unix(),
		}

		if err2 := models.InsertOneRoomBasic(&rb); err2 != nil {
			log.Printf("插入RoomBasic失败: %v", err2)
			ctx.JSON(http.StatusOK, response.Fail(nil, "添加好友失败"))
			return
		}

		// 2. 插入UserRoom
		ur := models.UserRoom{
			UserIdentity: userClaim.(*common.UserClaim).Identity,
			RoomIdentity: rb.Identity,
			RoomType:     1,
			CreateAt:     time.Now().Unix(),
			UpdateAt:     time.Now().Unix(),
		}

		if err = models.InsertOneUserRoom(&ur); err != nil {
			log.Printf("插入UserRoom失败: %v", err)
			ctx.JSON(http.StatusOK, response.Fail(nil, "添加好友失败"))
			return
		}

		ur = models.UserRoom{
			UserIdentity: ub.Identity,
			RoomIdentity: rb.Identity,
			RoomType:     1,
			CreateAt:     time.Now().Unix(),
			UpdateAt:     time.Now().Unix(),
		}

		if err = models.InsertOneUserRoom(&ur); err != nil {
			log.Printf("插入UserRoom失败: %v", err)
			ctx.JSON(http.StatusOK, response.Fail(nil, "添加好友失败"))
			return
		}

		ctx.JSON(http.StatusOK, response.Success(nil, "添加好友成功"))
	}
}

// 删除好友
func UserDelete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取参数
		username := ctx.Param("username")

		// 参数校验
		if username == "" {
			ctx.JSON(http.StatusOK, response.Fail(nil, "要查询的用户名能为空"))
			return
		}

		// 删除UserRoom
		userClaim, exists := ctx.Get("user_claim")
		if !exists {
			ctx.JSON(http.StatusOK, response.Fail(nil, "获取用户信息失败"))
			return
		}

		friend, err := models.GetUserBasicByUsername(username)
		if err != nil {
			ctx.JSON(http.StatusOK, response.Fail(nil, "要删除的用户不存在"))
			return
		}

		// 1. 找到自己和要删除用户共同的roomIdentity
		roomIdentity, err := models.GetRoomIdentityByUserIdentities(userClaim.(*common.UserClaim).Identity, friend.Identity)
		if err != nil || roomIdentity == "" {
			ctx.JSON(http.StatusOK, response.Fail(nil, "不是好友，无需删除"))
			return
		}

		// 2. 根据roomIdentity删除UserRoom
		if err2 := models.DeleteUserRoomByRoomIdentity(roomIdentity); err2 != nil {
			ctx.JSON(http.StatusOK, response.Fail(nil, "删除好友时发生错误,请重试"+err2.Error()))
			return
		}

		// 删除RoomBasic
		if err3 := models.DeleteRoomBasicByRoomIdentity(roomIdentity); err3 != nil {
			ctx.JSON(http.StatusOK, response.Fail(nil, "删除好友时发生错误,请重试"+err3.Error()))
			return
		}

		ctx.JSON(http.StatusOK, response.Success(nil, "删除好友成功"))
	}
}

// TODO:获取所有好友
func UserAllFirend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取用户identity
		// userIdentity := ctx.MustGet("user_claim").(*common.UserClaim).Identity

	}
}

// 创建群聊
func UserCreateChat() gin.HandlerFunc {
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
