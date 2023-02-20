package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prynnekey/gin-im/common"
)

func AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			c.Abort()
			return
		}
		claim, err := common.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			c.Abort()
			return
		}
		c.Set("user_claim", claim)
		c.Next()
	}
}
