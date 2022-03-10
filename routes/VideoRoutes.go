package routes

import (
	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/controller"
	"kuukaa.fun/danmu-v4/middleware"
)

func GetVideoRoutes(route *gin.RouterGroup) {
	//video 信息的增删改查接口
	video := route.Group("/video")
	{
		video.GET("/get", middleware.UidMiddleware(), controller.GetVideoByID)
		video.GET("/recommend/get", controller.GetRecommendVideo)
		video.GET("/list/get", controller.GetVideoList)
		video.GET("/user/get", controller.GetVideoListByUserID)
		videoAuth := video.Group("")
		videoAuth.Use(middleware.AuthMiddleware())
		{
			videoAuth.GET("/status", controller.GetVideoStatus)
			videoAuth.GET("/collect/get", controller.GetCollectVideo)
			videoAuth.GET("/upload/get", controller.GetMyUploadVideo)
			videoAuth.POST("/update", controller.ModifyVideoInfo) //只更新视频信息
			videoAuth.POST("/delete", controller.DeleteVideo)
			videoAuth.POST("/upload", controller.UploadVideoInfo)
		}
	}
}
