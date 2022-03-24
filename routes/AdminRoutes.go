package routes

import (
	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/controller/manage"
	"kuukaa.fun/danmu-v4/middleware"
	"kuukaa.fun/danmu-v4/util"
)

func GetAdminRoutes(route *gin.RouterGroup) {
	//管理员接口
	admin := route.Group("/admin")
	{
		admin.POST("/login", manage.AdminLogin)
		superAdminAuth := admin.Group("")
		superAdminAuth.Use(middleware.AdminMiddleware(util.SuperAdmin))
		{
			superAdminAuth.POST("/add", manage.AddAdmin) //添加管理员
			superAdminAuth.GET("/list", manage.GetAdminList)
			superAdminAuth.POST("/delete", manage.DeleteAdmin)
			superAdminAuth.POST("/user/delete", manage.AdminDeleteUser)
			superAdminAuth.POST("/partition/add", manage.AddPartition)
			superAdminAuth.POST("/partition/delete", manage.DeletePartition)
			superAdminAuth.GET("/update/get", manage.CheckUpdate) //检查更新
			config := superAdminAuth.Group("/config")
			{
				config.GET("/oss/get", manage.GetOssConfig)
				config.GET("/email/get", manage.GetEmailConfig)
				config.GET("/other/get", manage.GetOtherConfig)
				config.GET("/skin/get", manage.GetSkinInfoList)
				config.POST("/skin/upload", manage.UploadSkin) //上传主题
				config.POST("/skin/apply", manage.ApplySkin)   //应用主题
				config.POST("/skin/delete", manage.DeleteSkin) //删除主题
				config.POST("/oss/set", manage.SetOssConfig)
				config.POST("/email/set", manage.SetEmailConfig)
				config.POST("/admin/set", manage.SetAdminConfig)
				config.POST("/other/set", manage.SetOtherConfig)
			}
		}
		adminAuth := admin.Group("")
		adminAuth.Use(middleware.AdminMiddleware(util.Admin))
		{
			adminAuth.GET("/data", manage.GetRecentWebsiteData)
			adminAuth.GET("/data/total", manage.GetTotalWebsiteData)
			adminAuth.GET("/opinion/list", manage.GetOpinionList) //获取反馈列表
			adminAuth.GET("/user/search", manage.AdminSearchUser) //管理员搜索用户
			adminAuth.GET("/user/list", manage.GetUserList)
			adminAuth.POST("/user/modify", manage.AdminModifyUser)
			adminAuth.POST("/video/cover/upload", manage.AdminUploadCover) //管理员上传封面
			adminAuth.GET("/video/search", manage.AdminSearchVideo)        //管理员搜索视频
			adminAuth.GET("/video/list", manage.AdminGetVideoList)
			adminAuth.POST("/video/import", manage.ImportVideo)
			adminAuth.GET("/video/resource/list", manage.GetResourceList)
			adminAuth.POST("/video/resource/upload", manage.AdminUploadVideo) //管理员上传视频
			adminAuth.POST("/video/resource/delete", manage.DeleteResource)
			adminAuth.POST("/video/resource/import", manage.ImportResource)
			adminAuth.POST("/video/delete", manage.AdminDeleteVideo)
			adminAuth.POST("/announce/add", manage.AddAnnounce)
			adminAuth.POST("/announce/delete", manage.DeleteAnnounce)
			adminAuth.POST("/carousel/upload/img", manage.UploadCarousel)
			adminAuth.POST("/carousel/upload/info", manage.UploadCarouselInfo)
			adminAuth.POST("/carousel/delete", manage.DeleteCarousel)
			adminAuth.POST("/collection/delete", manage.AdminDeleteCollection)
		}

		auditorAuth := admin.Group("")
		auditorAuth.Use(middleware.AdminMiddleware(util.Auditor))
		{
			auditorAuth.GET("/review/list", manage.GetReviewVideoList)
			auditorAuth.POST("/review", manage.ReviewVideo)
			auditorAuth.GET("/announce/list", manage.AdminGetAnnounce)
			auditorAuth.GET("/carousel", manage.AdminGetCarousel)
			auditorAuth.GET("/video", manage.GetReviewVideoByID)
		}
	}
}
