package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/common"
	"kuukaa.fun/danmu-v4/model"
)

/*********************************************************
** 函数功能: 用户认证中间件
** 日    期:2021/6/7
**********************************************************/
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取Authorization，header
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "请先登录"})
			//抛弃请求
			ctx.Abort()
			return
		}

		tokenString = tokenString[7:]
		token, claims, err := common.ParseUserToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "请先登录"})
			//抛弃请求
			ctx.Abort()
			return
		}

		//通过token获取userID
		userId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)

		//验证用户是否存在
		if user.ID == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "请先登录"})
			//抛弃请求
			ctx.Abort()
			return
		}

		//如果用户存在，将用户信息写入上下文
		ctx.Set("id", user.ID)
		ctx.Set("user", user)
		ctx.Next()
	}
}
