package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/common"
	"kuukaa.fun/danmu-v4/model"
)

/*********************************************************
** 函数功能: 管理员认证中间件
** 日    期:2021/8/1
**********************************************************/
func AdminMiddleware(authority int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取Authorization，header
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			//抛弃请求
			ctx.Abort()
			return
		}

		tokenString = tokenString[7:]
		token, claims, err := common.ParseAdminToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			//抛弃请求
			ctx.Abort()
			return
		}

		//通过token获取管理员
		adminID := claims.AdminID
		if adminID == 0 {
			ctx.Next()
		} else {
			DB := common.GetDB()
			var admin model.Admin
			DB.First(&admin, adminID)
			//验证管理员是否存在
			if admin.ID == 0 {
				ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
				//抛弃请求
				ctx.Abort()
				return
			}
			//如果用户存在，比对目标权限
			if admin.Authority < authority {
				ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
				//抛弃请求
				ctx.Abort()
				return
			}
			ctx.Next()
		}
	}
}
