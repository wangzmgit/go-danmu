package routes

import (
	"wzm/danmu3.0/controller"
	"wzm/danmu3.0/middleware"

	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware())

	v1 := r.Group("/api/v1")
	{
		user := v1.Group("/user")
		{
			user.GET("/info/other", controller.GetUserInfoByID)
			user.POST("/register", controller.Register) //用户注册
			user.POST("/login", controller.Login)       //用户登录
			//需要用户登录
			userAuth := user.Group("")
			userAuth.Use(middleware.AuthMiddleware())
			{
				userAuth.GET("/info/get", controller.UserInfo) //用户获取个人信息
				userAuth.POST("/info/modify", controller.ModifyInfo)
			}
		}

		code := v1.Group("/code")
		{
			code.POST("/send", controller.SendCode)
		}

		//video 信息的增删改查接口
		video := v1.Group("/video")
		{
			video.GET("/get", controller.GetVideoByID)
			video.GET("/recommend/get", controller.GetRecommendVideo)
			video.GET("/list/get", controller.GetVideoList)
			video.GET("/user/get", controller.GetVideoListByUserID)
			videoAuth := video.Group("")
			videoAuth.Use(middleware.AuthMiddleware())
			{
				videoAuth.GET("/status", controller.GetVideoStatus)
				videoAuth.GET("/collect/get", controller.GetCollectVideo)
				videoAuth.GET("/upload/get", controller.GetMyUploadVideo)
				videoAuth.POST("/update/request", controller.UpdateRequest)
				videoAuth.POST("/update", controller.ModifyVideoInfo) //只更新视频信息
				videoAuth.POST("/delete", controller.DeleteVideo)
				videoAuth.POST("/upload", controller.UploadVideoInfo)
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

		//点赞收藏
		interactive := v1.Group("/interactive")
		interactive.Use(middleware.AuthMiddleware())
		{
			interactive.GET("/video", controller.GetVideoInteractiveData) //获取点赞收藏关注的交互数据
			interactive.POST("/collect/add", controller.Collect)
			interactive.POST("/collect/cancel", controller.CancelCollect)
			interactive.POST("/like/add", controller.Like)
			interactive.POST("/like/cancel", controller.Dislike)
		}

		//关注粉丝
		v1.GET("follow/following", controller.GetFollowingByID) //关注列表
		v1.GET("follow/followers", controller.GetFollowersByID) //粉丝列表
		v1.GET("follow/count", controller.GetFollowCount)
		follow := v1.Group("/follow")
		follow.Use(middleware.AuthMiddleware())
		{
			follow.GET("/status", controller.GetFollowStatus)
			follow.POST("", controller.Following)
			follow.POST("/cancel", controller.UnFollow)
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

		message := v1.Group("/message")
		message.Use(middleware.AuthMiddleware())
		{
			message.GET("/announce", controller.GetAnnounce)
			message.GET("/list", controller.GetMessageList)
			message.GET("/details", controller.GetMessageDetails)
			message.POST("/send", controller.SendMessage)
		}
		danmaku := v1.Group("/danmaku")
		{
			danmaku.GET("/get", controller.GetDanmaku)
			danmaku.POST("send", middleware.AuthMiddleware(), controller.SendDanmaku)
		}
		//其他接口
		v1.GET("search", controller.Search)
		v1.GET("carousel", controller.GetCarousel)

		//管理员接口
		admin := v1.Group("/admin")
		{
			admin.POST("/login", controller.AdminLogin)
			superAdminAuth := admin.Group("")
			superAdminAuth.Use(middleware.AdminMiddleware(controller.SuperAdmin))
			{
				admin.POST("/add", controller.AddAdmin) //添加管理员
				admin.GET("/list", controller.GetAdminList)
				admin.POST("/delete", controller.DeleteAdmin)
				admin.POST("/user/delete", controller.AdminDeleteUser)
			}
			adminAuth := admin.Group("")
			adminAuth.Use(middleware.AdminMiddleware(controller.Admin))
			{
				admin.GET("/user/list", controller.GetUserList)
				admin.POST("/user/modify", controller.AdminModifyUser)
				admin.GET("/video/list", controller.AdminGetVideoList)
				admin.POST("/video/delete", controller.AdminDeleteVideo)
				admin.POST("/announce/add", controller.AddAnnounce)
				admin.POST("/announce/delete", controller.DeleteAnnounce)
				admin.POST("/carousel/upload/img", controller.UploadCarousel)
				admin.POST("/carousel/upload/info", controller.UploadCarouselInfo)
				admin.POST("/carousel/delete", controller.DeleteCarousel)
			}

			auditorAuth := admin.Group("")
			auditorAuth.Use(middleware.AdminMiddleware(controller.Auditor))
			{
				admin.GET("/review/list", controller.GetReviewVideoList)
				admin.POST("/review", controller.ReviewVideo)
				admin.GET("/announce/list", controller.AdminGetAnnounce)
				admin.GET("/carousel", controller.AdminGetCarousel)
			}
		}
	}
	return r
}
