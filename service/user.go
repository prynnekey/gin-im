package service

import (
	"net/http"

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
