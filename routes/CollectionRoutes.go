package routes

import (
	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/controller"
	"kuukaa.fun/danmu-v4/middleware"
)

func GetCollectionRoutes(route *gin.RouterGroup) {
	//合集
	collection := route.Group("/collection")
	{
		collection.GET("/get", controller.GetCollectionByID)
		collection.GET("/list/get", controller.GetCollectionList)
		collection.GET("/video/get", controller.GetCollectionContent)
		collectionAuth := collection.Group("")
		collectionAuth.Use(middleware.AuthMiddleware())
		{
			collectionAuth.POST("/modify", controller.ModifyCollection)
			collectionAuth.POST("/delete", controller.DeleteCollection)
			collectionAuth.POST("/create", controller.CreateCollection)
			collectionAuth.POST("/video/add", controller.AddVideoToCollection)     //添加视频
			collectionAuth.POST("/video/delete", controller.DeleteCollectionVideo) //删除视频
			collectionAuth.GET("/video/add/list", controller.GetCanAddVideo)       //可以添加的视频
			collectionAuth.GET("/create/list", controller.GetCreateCollectionList) //创建的合集列表
		}
	}
}
