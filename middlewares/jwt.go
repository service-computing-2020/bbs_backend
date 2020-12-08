package middlewares

import (
	"github.com/service-computing-2020/bbs_backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func VerifyJWT() gin.HandlerFunc{
	return func(c *gin.Context) {
		isError := false
		var data interface{}
		token := c.Query("token")
		if token == "" {
			isError = true
			c.JSON(http.StatusBadRequest, gin.H {
				"code": 400,
				"msg": "缺少token",
				"data":data,
			})
		} else {
			claims, err := service.ParseToken(token)
			if err != nil {
				isError = true
				c.JSON(http.StatusInternalServerError, gin.H {
					"code":500,
					"msg":"token校验发生错误",
					"data":data,
				})
			} else if time.Now().Unix() > claims.ExpiresAt {
				isError = true
				c.JSON(http.StatusForbidden, gin.H{
					"code":403,
					"msg":"token已过期",
					"data":data,
				})
			}
		}

		if isError {
			c.Abort()
			return
		}
		c.Next()
	}
}