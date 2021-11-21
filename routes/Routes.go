package routes

import (
	"wzm/danmu3.0/controller"
	admin_controller "wzm/danmu3.0/controller/admin"
	"wzm/danmu3.0/middleware"
	"wzm/danmu3.0/util"

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
				userAuth.POST("/password/modify", controller.ModifyPassword)
			}
		}

		code := v1.Group("/code")
		{
			code.POST("/send", controller.SendCode)
			code.POST("/send/myself", middleware.AuthMiddleware(), controller.SendCodeToMyself)
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
			danmaku.POST("/send", middleware.AuthMiddleware(), controller.SendDanmaku)
		}

		collection := v1.Group("/collection")
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

		//其他接口
		v1.GET("search", controller.Search)
		v1.GET("carousel", controller.GetCarousel)

		//管理员接口
		admin := v1.Group("/admin")
		{
			admin.POST("/login", admin_controller.AdminLogin)
			superAdminAuth := admin.Group("")
			superAdminAuth.Use(middleware.AdminMiddleware(util.SuperAdmin))
			{
				superAdminAuth.POST("/add", admin_controller.AddAdmin) //添加管理员
				superAdminAuth.GET("/list", admin_controller.GetAdminList)
				superAdminAuth.POST("/delete", admin_controller.DeleteAdmin)
				superAdminAuth.POST("/user/delete", admin_controller.AdminDeleteUser)
			}
			adminAuth := admin.Group("")
			adminAuth.Use(middleware.AdminMiddleware(util.Admin))
			{
				adminAuth.GET("/user/list", admin_controller.GetUserList)
				adminAuth.POST("/user/modify", admin_controller.AdminModifyUser)
				adminAuth.GET("/video/list", admin_controller.AdminGetVideoList)
				adminAuth.POST("/video/add", admin_controller.ImportVideo)
				adminAuth.POST("/video/delete", admin_controller.AdminDeleteVideo)
				adminAuth.POST("/announce/add", admin_controller.AddAnnounce)
				adminAuth.POST("/announce/delete", admin_controller.DeleteAnnounce)
				adminAuth.POST("/carousel/upload/img", admin_controller.UploadCarousel)
				adminAuth.POST("/carousel/upload/info", admin_controller.UploadCarouselInfo)
				adminAuth.POST("/carousel/delete", admin_controller.DeleteCarousel)
			}

			auditorAuth := admin.Group("")
			auditorAuth.Use(middleware.AdminMiddleware(util.Auditor))
			{
				auditorAuth.GET("/review/list", admin_controller.GetReviewVideoList)
				auditorAuth.POST("/review", admin_controller.ReviewVideo)
				auditorAuth.GET("/announce/list", admin_controller.AdminGetAnnounce)
				auditorAuth.GET("/carousel", admin_controller.AdminGetCarousel)
			}
		}

		v2 := r.Group("/api/v2")
		{
			//评论回复
			v2.GET("/comment/get", controller.GetCommentsV2)
			v2.GET("/comment/reply", controller.GetReplyDetailsV2)
		}
	}
	return r
}
