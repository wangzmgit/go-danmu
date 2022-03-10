package routes

import (
	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/controller"
	"kuukaa.fun/danmu-v4/middleware"
)

func GetUserRoutes(route *gin.RouterGroup) {
	user := route.Group("/user")
	{
		user.GET("/info/other", controller.GetUserInfoByID)
		user.POST("/register", controller.Register)      //用户注册
		user.POST("/login", controller.Login)            //用户登录
		user.POST("/email/login", controller.EmailLogin) //邮箱登录
		//需要用户登录
		userAuth := user.Group("")
		userAuth.Use(middleware.AuthMiddleware())
		{
			userAuth.GET("/info/get", controller.UserInfo) //用户获取个人信息
			userAuth.POST("/info/modify", controller.ModifyInfo)
			userAuth.POST("/password/modify", controller.ModifyPassword)
		}
	}
}
