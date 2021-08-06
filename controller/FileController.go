package controller

import (
	"os"
	"path"
	"strconv"
	"time"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/util"

	"github.com/gin-gonic/gin"
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
	// 拼接上传图片的路径信息
	localFileName := "./file/avatar/" + avatar.Filename
	objectName := "avatar/" + avatar.Filename
	success, url := util.UploadOSS(localFileName, objectName)
	if success {
		id, _ := ctx.Get("id")
		DB := common.GetDB()
		DB.Model(model.User{}).Where("id = ?", id).Update("avatar", url)
		response.Success(ctx, nil, "ok")
	} else {
		response.Fail(ctx, nil, "上传失败")
	}
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
	// 拼接上传图片的路径信息
	localFileName := "./file/cover/" + cover.Filename
	objectName := "cover/" + cover.Filename
	success, url := util.UploadOSS(localFileName, objectName)
	if success {
		response.Success(ctx, gin.H{"url": url}, "ok")
	} else {
		response.Fail(ctx, nil, "上传失败")
	}
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
	// 拼接上传图片的路径信息
	localFileName := "./file/video/" + video.Filename
	objectName := "video/" + video.Filename
	success, url := util.UploadOSS(localFileName, objectName)
	if success {
		uid, _ := ctx.Get("id")
		var video model.Video
		DB := common.GetDB()
		DB.Where("id = ?", vid).Last(&video)
		if video.ID == 0 || video.Uid != uid {
			response.Fail(ctx, nil, "视频不存在")
			return
		}
		//开始事务
		tx := DB.Begin()
		if err := tx.Model(&video).Update("video", url).Error; err != nil {
			util.Logfile("[Error]", " upload video error "+err.Error())
			tx.Rollback()
			response.Fail(ctx, nil, "上传失败")
			return
		}
		//创建新的审核状态
		if err := tx.Model(&model.Review{}).Where("vid = ?", vid).Updates(map[string]interface{}{"status": 1000}).Error; err != nil {
			tx.Rollback()
			response.Fail(ctx, nil, "上传失败")
			return
		}
		tx.Commit()
		response.Success(ctx, nil, "ok")
	} else {
		response.Fail(ctx, nil, "上传失败")
	}
}
