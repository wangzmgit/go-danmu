package controller

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
	util.Logfile("[Info]", " User "+strconv.Itoa(int(uid.(uint)))+" | "+ctx.ClientIP()+" | "+objectName)

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

	uid, _ := ctx.Get("id")
	localFileName := "./file/cover/" + cover.Filename
	objectName := "cover/" + cover.Filename
	//记录日志
	util.Logfile("[Info]", " User "+strconv.Itoa(int(uid.(uint)))+" | "+ctx.ClientIP()+" | "+objectName)

	res := service.UploadCoverService(localFileName, objectName)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 上传视频
** 日    期:2021/7/14
**********************************************************/
func UploadVideo(ctx *gin.Context) {
	urls := map[string]string{
		"res360":   "",
		"res480":   "",
		"res720":   "",
		"res1080":  "",
		"original": "",
	}
	vid, _ := strconv.Atoi(ctx.PostForm("vid"))
	if vid <= 0 {
		response.Fail(ctx, nil, "VID格式有误")
		return
	}
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

	uid, _ := ctx.Get("id")
	// 拼接上传图片的路径信息
	localFileName := "./file/video/" + video.Filename
	objectName := "video/" + video.Filename
	//获取url
	if viper.GetString("transcoding.coding") == "hls" {
		if viper.GetBool("aliyunoss.storage") {
			urls["original"] = service.GetUrl() + "video/" + videoName + "/" + "index.m3u8"
		} else {
			urls["original"] = service.GetUrl() + "output/" + videoName + "/" + "index.m3u8"
		}
	} else {
		urls["original"] = service.GetUrl() + objectName
	}
	//记录日志
	util.Logfile("[Info]", " User "+strconv.Itoa(int(uid.(uint)))+" | "+ctx.ClientIP()+" | "+objectName)
	res := service.UploadVideoService(urls, vid, uid.(uint))
	//启动转码服务或上传服务
	if viper.GetString("transcoding.coding") == "hls" {
		go service.Transcoding(video.Filename, vid)
	} else {
		if viper.GetBool("aliyunoss.storage") {
			go service.UploadVideoToOSS(localFileName, objectName, vid)
		} else {
			service.CompleteUpload(vid)
		}
	}

	response.HandleResponse(ctx, res)
}
