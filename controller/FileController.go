package controller

import (
	"os"
	"path"
	"strconv"
	"time"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/service"
	"wzm/danmu3.0/util"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

/*********************************************************
** 函数功能: 上传头像
** 日    期:2021/7/13
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
	// 拼接上传图片的路径信息
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
	video.Filename = util.RandomString(3) + strconv.FormatInt(time.Now().UnixNano(), 10) + suffix
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

	var url string
	uid, _ := ctx.Get("id")
	// 拼接上传图片的路径信息
	localFileName := "./file/video/" + video.Filename
	objectName := "video/" + video.Filename

	//启用hls
	if viper.GetString("server.coding") == "hls" {
		url = service.Transcoding(video.Filename, vid)
	} else {
		go service.UploadVideoToOSS(localFileName, objectName, vid)
		if len(viper.GetString("aliyunoss.domain")) == 0 {
			url = "http://" + viper.GetString("aliyunoss.bucket") + "." + viper.GetString("aliyunoss.endpoint") + "/" + objectName
		} else {
			url = "http://" + viper.GetString("aliyunoss.domain") + "/" + objectName
		}
	}

	//记录日志
	util.Logfile("[Info]", " User "+strconv.Itoa(int(uid.(uint)))+" | "+ctx.ClientIP()+" | "+objectName)

	res := service.UploadVideoService(url, vid, uid.(uint))
	response.HandleResponse(ctx, res)
}
