package service

import (
	"net/http"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/util"

	"github.com/gin-gonic/gin"
)

/*********************************************************
** 函数功能: 上传头像
** 日    期: 2021年11月11日17:31:42
**********************************************************/
func UploadAvatarService(localFileName string, objectName string, uid uint) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}

	success, url := UploadOSS(localFileName, objectName)
	if success {
		DB := common.GetDB()
		DB.Model(model.User{}).Where("id = ?", uid).Update("avatar", url)
		return res
	}
	res.HttpStatus = http.StatusBadRequest
	res.Code = response.FailCode
	return res
}

/*********************************************************
** 函数功能: 上传封面
** 日    期: 2021年11月11日17:38:07
**********************************************************/
func UploadCoverService(localFileName string, objectName string) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}
	success, url := UploadOSS(localFileName, objectName)

	if success {
		res.Data = gin.H{"url": url}
		return res
	}
	res.HttpStatus = http.StatusBadRequest
	res.Code = response.FailCode
	return res
}

/*********************************************************
** 函数功能: 上传封面
** 日    期: 2021年11月11日17:38:07
**********************************************************/
func UploadVideoService(url string, vid int, uid uint) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}

	var videoInfo model.Video
	DB := common.GetDB()
	DB.Where("id = ?", vid).First(&videoInfo)
	if videoInfo.ID == 0 || videoInfo.Uid != uid {
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = "视频不存在"
		return res
	}
	//开始事务
	tx := DB.Begin()
	if err := tx.Model(&videoInfo).Update("video", url).Error; err != nil {
		util.Logfile("[Error]", " upload video error "+err.Error())
		tx.Rollback()
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = "上传失败"
		return res
	}
	//创建新的审核状态
	if err := tx.Model(&model.Review{}).Where("vid = ?", vid).Updates(map[string]interface{}{"status": 800}).Error; err != nil {
		tx.Rollback()
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = "上传失败"
		return res
	}
	tx.Commit()
	return res
}

/*********************************************************
** 函数功能: 上传轮播图
** 日    期: 2021年11月12日12:24:55
**********************************************************/
func UploadCarouselService(localFileName string, objectName string) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}

	success, url := UploadOSS(localFileName, objectName)
	if success {
		res.Data = gin.H{"url": url}
		return res
	}
	res.HttpStatus = http.StatusBadRequest
	res.Code = response.FailCode
	res.Msg = "上传失败"
	return res
}

/*********************************************************
** 函数功能: 完成视频上传
** 日    期:2021/9/16
**********************************************************/
func CompleteUpload(vid int) {
	DB := common.GetDB()
	DB.Model(&model.Review{}).Where("vid = ?", vid).Updates(map[string]interface{}{"status": 1000})
}
