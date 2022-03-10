package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/controller"
	"kuukaa.fun/danmu-v4/middleware"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware())
	v1 := r.Group("/api/v1")
	{
		GetUserRoutes(v1)       //获取用户路由
		GetVideoRoutes(v1)      //获取视频路由
		GetAdminRoutes(v1)      // 获取管理员路由
		GetCollectionRoutes(v1) //获取合集路由

		code := v1.Group("/code")
		{
			code.POST("/send", controller.SendCode)
			code.POST("/send/login", controller.SendLoginCode)
			code.POST("/send/myself", middleware.AuthMiddleware(), controller.SendCodeToMyself)
		}

		//点赞收藏
		interactive := v1.Group("/interactive")
		interactive.Use(middleware.AuthMiddleware())
		{
			interactive.POST("/collect/add", controller.Collect)
			interactive.POST("/collect/cancel", controller.CancelCollect)
			interactive.POST("/like/add", controller.Like)
			interactive.POST("/like/cancel", controller.Dislike)
		}

		//关注粉丝
		follow := v1.Group("/follow")
		{
			follow.GET("/following", controller.GetFollowingByID) //关注列表
			follow.GET("/followers", controller.GetFollowersByID) //粉丝列表
			follow.GET("/count", controller.GetFollowCount)
			followAuth := follow.Group("")
			followAuth.Use(middleware.AuthMiddleware())
			{
				followAuth.GET("/status", controller.GetFollowStatus)
				followAuth.POST("", controller.Following)
				followAuth.POST("/cancel", controller.UnFollow)
			}
		}

		//评论回复
		v1.GET("comment/get", controller.GetComments)
		comment := v1.Group("/comment")
		comment.Use(middleware.AuthMiddleware())
		{
			comment.POST("", controller.Comment) //评论
			comment.POST("/reply", controller.Reply)
			comment.POST("/delete", controller.DeleteComment)
			comment.POST("/reply/delete", controller.DeleteReply)
		}

		// 消息
		message := v1.Group("/message")
		message.GET("/ws", controller.GetMsgConnect)
		message.Use(middleware.AuthMiddleware())
		{
			message.GET("/announce", controller.GetAnnounce)
			message.GET("/list", controller.GetMessageList)
			message.GET("/details", controller.GetMessageDetails)
			message.POST("/send", controller.SendMessage)
			message.POST("/read", controller.ReadMessageService)
		}

		// 弹幕
		danmaku := v1.Group("/danmaku")
		{
			danmaku.GET("/get", controller.GetDanmaku)
			danmaku.POST("/send", middleware.AuthMiddleware(), controller.SendDanmaku)
		}

		//其他接口
		v1.GET("/search", controller.Search)
		v1.GET("/carousel", controller.GetCarousel)
		v1.GET("/partition/list", controller.GetPartitionList)
		v1.GET("/partition/all", controller.GetAllPartition)
		v1.POST("opinion", controller.CreateOpinion)
		v1.POST("opinion/site", middleware.AuthMiddleware(), controller.CreateOpinionOnSite)

		v2 := r.Group("/api/v2")
		{
			//评论回复
			v2.GET("/comment/get", controller.GetCommentsV2)
			v2.GET("/comment/reply", controller.GetReplyDetailsV2)

			messageV2 := v2.Group("/message")
			messageV2.Use(middleware.AuthMiddleware())
			{
				messageV2.GET("/details", controller.GetMessageDetailsV2) //v2获取私信内容
			}
		}

		//文件上传相关接口
		file := v1.Group("/file")
		file.Use(middleware.AuthMiddleware())
		{
			file.POST("/avatar", controller.UploadAvatar)
			file.POST("/cover", controller.UploadCover)
			file.POST("/video", controller.UploadVideo)
		}
		//获取静态文件
		r.StaticFS("/api/avatar", http.Dir("./file/avatar"))
		r.StaticFS("/api/cover", http.Dir("./file/cover"))
		r.StaticFS("/api/video", http.Dir("./file/video"))
		r.StaticFS("/api/carousel", http.Dir("./file/carousel"))
		r.StaticFS("/api/output", http.Dir("./file/output"))
	}
	return r
}
