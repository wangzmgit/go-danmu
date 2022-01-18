package manage

import (
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/service"
	"kuukaa.fun/danmu-v4/util"
)

/*********************************************************
** 函数功能: 上传轮播图
** 日    期: 2021/8/4
**********************************************************/
func UploadCarousel(ctx *gin.Context) {
	carousel, err := ctx.FormFile("carousel")
	if err != nil {
		response.Fail(ctx, nil, "图片上传失败")
		return
	}
	suffix := path.Ext(carousel.Filename)
	if suffix != ".jpg" && suffix != ".jpeg" && suffix != ".png" {
		response.CheckFail(ctx, nil, "图片不符合要求")
		return
	}
	carousel.Filename = util.RandomString(3) + strconv.FormatInt(time.Now().UnixNano(), 10) + suffix
	errSave := ctx.SaveUploadedFile(carousel, "./file/carousel/"+carousel.Filename)
	if errSave != nil {
		response.Fail(ctx, nil, "图片保存失败")
		return
	}
	fileInfo, err := os.Stat("./file/carousel/" + carousel.Filename)
	//大小限制到5M
	if fileInfo == nil || fileInfo.Size() > 1024*1024*5 || err != nil {
		response.CheckFail(ctx, nil, "图片大小不符合要求")
		return
	}
	// 拼接上传图片的路径信息
	localFileName := "./file/carousel/" + carousel.Filename
	objectName := "carousel/" + carousel.Filename

	res := service.UploadCarouselService(localFileName, objectName)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 管理员上传视频
** 日    期: 2022年1月14日13:41:22
**********************************************************/
func AdminUploadVideo(ctx *gin.Context) {
	video, err := ctx.FormFile("video")
	if err != nil {
		response.Fail(ctx, nil, "视频上传失败")
		return
	}
	suffix := path.Ext(video.Filename)
	if suffix != ".mp4" {
		response.CheckFail(ctx, nil, "请上传mp4格式文件")
		return
	}
	//仅文件名(不含后缀)
	videoName := util.RandomString(3) + strconv.FormatInt(time.Now().UnixNano(), 10)
	video.Filename = videoName + suffix
	errSave := ctx.SaveUploadedFile(video, "./file/video/"+video.Filename)
	if errSave != nil {
		response.Fail(ctx, nil, "视频保存失败")
		return
	}
	fileInfo, err := os.Stat("./file/video/" + video.Filename)
	//大小限制到500M
	if fileInfo == nil || fileInfo.Size() > 1024*1024*500 || err != nil {
		response.CheckFail(ctx, nil, "视频大小不符合要求")
		return
	}

	// 拼接上传图片的路径信息
	localFileName := "./file/video/" + video.Filename
	objectName := "video/" + video.Filename
	//获取url
	url := service.GetUrl() + objectName
	//记录日志
	util.Logfile("[Info]", " admin upload video"+" | "+ctx.ClientIP()+" | "+objectName)
	//启动上传服务
	if viper.GetBool("aliyunoss.storage") {
		go service.UploadOSS(localFileName, objectName)
	}

	response.Success(ctx, gin.H{"url": url}, "ok")
}

/*********************************************************
** 函数功能: 上传视频封面图
** 日    期: 2022年1月14日14:01:45
**********************************************************/
func AdminUploadCover(ctx *gin.Context) {
	cover, err := ctx.FormFile("cover")
	if err != nil {
		response.Fail(ctx, nil, "图片上传失败")
		return
	}
	suffix := path.Ext(cover.Filename)
	if suffix != ".jpg" && suffix != ".jpeg" && suffix != ".png" {
		response.CheckFail(ctx, nil, "图片不符合要求")
		return
	}
	//储存文件
	cover.Filename = util.RandomString(3) + strconv.FormatInt(time.Now().UnixNano(), 10) + suffix
	errSave := ctx.SaveUploadedFile(cover, "./file/cover/"+cover.Filename)
	if errSave != nil {
		response.Fail(ctx, nil, "图片保存失败")
		return
	}
	fileInfo, err := os.Stat("./file/cover/" + cover.Filename)
	//大小限制到5M
	if fileInfo == nil || fileInfo.Size() > 1024*1024*5 || err != nil {
		response.CheckFail(ctx, nil, "图片不符合要求")
		return
	}

	localFileName := "./file/cover/" + cover.Filename
	objectName := "cover/" + cover.Filename
	//记录日志
	util.Logfile("[Info]", " Admin upload cover "+" | "+ctx.ClientIP()+" | "+objectName)

	res := service.UploadCoverService(localFileName, objectName)
	response.HandleResponse(ctx, res)
}
