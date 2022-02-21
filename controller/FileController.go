package controller

import (
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"kuukaa.fun/danmu-v4/dto"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/service"
	"kuukaa.fun/danmu-v4/util"
)

/*********************************************************
** 函数功能: 上传头像
** 日    期: 2021/7/13
**********************************************************/
func UploadAvatar(ctx *gin.Context) {
	avatar, err := ctx.FormFile("avatar")
	if err != nil {
		response.Fail(ctx, nil, "图片上传失败")
		return
	}
	suffix := path.Ext(avatar.Filename)
	if suffix != ".jpg" && suffix != ".jpeg" && suffix != ".png" {
		response.CheckFail(ctx, nil, "图片不符合要求")
		return
	}
	avatar.Filename = util.RandomString(3) + strconv.FormatInt(time.Now().UnixNano(), 10) + suffix
	errSave := ctx.SaveUploadedFile(avatar, "./file/avatar/"+avatar.Filename)
	if errSave != nil {
		response.Fail(ctx, nil, "图片保存失败")
		return
	}
	fileInfo, err := os.Stat("./file/avatar/" + avatar.Filename)
	//大小限制到5M
	if fileInfo == nil || fileInfo.Size() > 1024*1024*5 || err != nil {
		response.CheckFail(ctx, nil, "图片不符合要求")
		return
	}

	uid, _ := ctx.Get("id")
	// 拼接上传图片的路径信息
	localFileName := "./file/avatar/" + avatar.Filename
	objectName := "avatar/" + avatar.Filename

	//记录日志
	util.Logfile(util.InfoLog, " User "+strconv.Itoa(int(uid.(uint)))+" | "+ctx.ClientIP()+" | "+objectName)

	res := service.UploadAvatarService(localFileName, objectName, uid.(uint))
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 上传视频封面图
** 日    期:2021/7/14
**********************************************************/
func UploadCover(ctx *gin.Context) {
	cover, err := ctx.FormFile("cover")
	if err != nil {
		response.Fail(ctx, nil, response.FileUploadFail)
		return
	}
	suffix := path.Ext(cover.Filename)
	if suffix != ".jpg" && suffix != ".jpeg" && suffix != ".png" {
		response.CheckFail(ctx, nil, response.FileCheckFail)
		return
	}
	//储存文件
	cover.Filename = util.RandomString(3) + strconv.FormatInt(time.Now().UnixNano(), 10) + suffix
	errSave := ctx.SaveUploadedFile(cover, "./file/cover/"+cover.Filename)
	if errSave != nil {
		response.Fail(ctx, nil, response.FileSaveFail)
		return
	}
	fileInfo, err := os.Stat("./file/cover/" + cover.Filename)
	//大小限制到5M
	if fileInfo == nil || fileInfo.Size() > 1024*1024*5 || err != nil {
		response.CheckFail(ctx, nil, response.FileSizeCheckFail)
		return
	}

	uid, _ := ctx.Get("id")
	localFileName := "./file/cover/" + cover.Filename
	objectName := "cover/" + cover.Filename
	//记录日志
	util.Logfile(util.InfoLog, " User "+strconv.Itoa(int(uid.(uint)))+" | "+ctx.ClientIP()+" | "+objectName)

	res := service.UploadCoverService(localFileName, objectName)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 上传视频
** 日    期:2021/7/14
**********************************************************/
func UploadVideo(ctx *gin.Context) {
	maxRes := 0 //最大分辨率
	vid, _ := strconv.Atoi(ctx.PostForm("vid"))
	if vid <= 0 {
		response.Fail(ctx, nil, response.ParameterError)
		return
	}
	video, err := ctx.FormFile("video")
	if err != nil {
		response.Fail(ctx, nil, response.FileUploadFail)
		return
	}
	suffix := path.Ext(video.Filename)
	if suffix != ".mp4" {
		response.CheckFail(ctx, nil, response.FileCheckFail)
		return
	}
	//仅文件名(不含后缀)
	videoName := util.RandomString(3) + strconv.FormatInt(time.Now().UnixNano(), 10)
	video.Filename = videoName + suffix
	errSave := ctx.SaveUploadedFile(video, "./file/video/"+video.Filename)
	if errSave != nil {
		response.Fail(ctx, nil, response.FileSaveFail)
		return
	}
	fileInfo, err := os.Stat("./file/video/" + video.Filename)
	//大小限制到500M
	if fileInfo == nil || fileInfo.Size() > 1024*1024*500 || err != nil {
		response.CheckFail(ctx, nil, response.FileSizeCheckFail)
		return
	}

	uid, _ := ctx.Get("id")
	// 拼接上传图片的路径信息
	localFileName := "./file/video/" + video.Filename
	objectName := "video/" + video.Filename
	urls, maxRes := GetUploadVideoUrls(videoName, localFileName, objectName, vid)
	//记录日志
	util.Logfile(util.InfoLog, " User "+strconv.Itoa(int(uid.(uint)))+" | "+ctx.ClientIP()+" | "+objectName)
	res := service.UploadVideoService(urls, vid, uid.(uint))
	//启动转码服务或上传服务
	if viper.GetString("transcoding.coding") == "hls" {
		go service.Transcoding(video.Filename, vid, maxRes)
	} else {
		if viper.GetBool("aliyunoss.storage") {
			go service.UploadVideoToOSS(localFileName, objectName, vid)
		} else {
			service.CompleteUpload(vid)
		}
	}

	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 获取上传视频文件的url
** 日    期: 2022年2月16日17:11:08
**********************************************************/
func GetUploadVideoUrls(videoName, localFileName, objectName string, vid int) (dto.ResDto, int) {
	var maxRes int
	var urls dto.ResDto
	if viper.GetString("transcoding.coding") == "hls" {
		if viper.GetBool("aliyunoss.storage") {
			if viper.GetInt("transcoding.max_res") == 0 {
				urls.Original = service.GetUrl() + "video/" + videoName + "/" + "index.m3u8"
			} else {
				urls, maxRes = service.GetUrlDifferentRes(videoName, localFileName, vid, true)
			}
		} else {
			if viper.GetInt("transcoding.max_res") == 0 {
				urls.Original = service.GetUrl() + "output/" + videoName + "/" + "index.m3u8"
			} else {
				urls, maxRes = service.GetUrlDifferentRes(videoName, localFileName, vid, false)
			}
		}
	} else {
		urls.Original = service.GetUrl() + objectName
	}
	return urls, maxRes
}
