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
		token, err := common.GenerateToken(ub.Identidy, ub.Username)
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
		count, err := models.GetUserBasicByUsername(username)
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
			Identidy:  common.GenerateUUID(),
			Username:  username,
			Password:  password,
			Create_at: time.Now().Unix(),
			Update_at: time.Now().Unix(),
		}
		err = models.InsertOneUserBasic(&ub)
		if err != nil {
			log.Printf("插入数据失败: %v", err)
			ctx.JSON(http.StatusOK, response.Fail(nil, "发生错误"+err.Error()))
			return
		}

		// 注册成功

		// 生成token
		token, err := common.GenerateToken(ub.Identidy, ub.Username)
		if err != nil {
			ctx.JSON(http.StatusOK, response.Fail(nil, "生成token时失败"+err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, response.Success(gin.H{"token": token}, "注册成功"))
	}
}
