package service

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"kuukaa.fun/danmu-v4/common"
	"kuukaa.fun/danmu-v4/model"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/util"
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

	if viper.GetBool("aliyunoss.storage") {
		success := UploadOSS(localFileName, objectName)
		if !success {
			res.HttpStatus = http.StatusBadRequest
			res.Code = response.FailCode
			return res
		}
	}

	url := GetUrl() + objectName
	DB := common.GetDB()
	DB.Model(model.User{}).Where("id = ?", uid).Update("avatar", url)
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

	if viper.GetBool("aliyunoss.storage") {
		success := UploadOSS(localFileName, objectName)
		if !success {
			res.HttpStatus = http.StatusBadRequest
			res.Code = response.FailCode
			return res
		}
	}

	res.Data = gin.H{"url": GetUrl() + objectName}
	return res
}

/*********************************************************
** 函数功能: 上传视频
** 日    期: 2021年11月11日17:38:07
**********************************************************/
func UploadVideoService(urls map[string]string, vid int, uid uint) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}
	var err error
	var videoInfo model.Video

	DB := common.GetDB()
	DB.Where("id = ?", vid).First(&videoInfo)
	if videoInfo.ID == 0 || videoInfo.Uid != uid {
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = "视频不存在"
		return res
	}
	//当前版本不支持普通用户上传分P，在上传前将该视频的其他视频资源删除
	DB.Where("vid = ?", vid).Delete(model.Resource{})
	//开始事务
	tx := DB.Begin()
	var newResource model.Resource
	if viper.GetString("transcoding.coding") == "hls" {
		newResource.Vid = uint(vid)
		newResource.Res360 = urls["res360"]
		newResource.Res480 = urls["res480"]
		newResource.Res720 = urls["res720"]
		newResource.Res1080 = urls["res1080"]
		newResource.Original = urls["original"]
	} else {
		//视频类型为mp4,不进行转码，分辨率为原始分辨率
		newResource.Vid = uint(vid)
		newResource.Original = urls["original"]
	}
	if err = tx.Model(&model.Resource{}).Create(&newResource).Error; err != nil {
		util.Logfile("[Error]", " upload video error "+err.Error())
		tx.Rollback()
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = "上传失败"
		return res
	}
	//创建新的审核状态
	if err = tx.Model(&model.Review{}).Where("vid = ?", vid).Updates(map[string]interface{}{"status": 800}).Error; err != nil {
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

	if viper.GetBool("aliyunoss.storage") {
		success := UploadOSS(localFileName, objectName)
		if !success {
			res.HttpStatus = http.StatusBadRequest
			res.Code = response.FailCode
			res.Msg = "上传失败"
			return res
		}
	}

	res.Data = gin.H{"url": GetUrl() + objectName}
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

/*********************************************************
** 函数功能: 获取文件的URL
** 日    期: 2022年1月5日16:49:02
**********************************************************/
func GetUrl() string {
	if viper.GetBool("aliyunoss.storage") {
		if len(viper.GetString("aliyunoss.domain")) == 0 {
			return "http://" + viper.GetString("aliyunoss.bucket") + "." + viper.GetString("aliyunoss.endpoint") + "/"
		} else {
			return "http://" + viper.GetString("aliyunoss.domain") + "/"
		}
	} else {
		if len(viper.GetString("aliyunoss.domain")) == 0 {
			return "/api/"
		} else {
			return "http://" + viper.GetString("aliyunoss.domain") + "/api/"
		}
	}
}

/*********************************************************
** 函数功能: 获取不同分辨率URL
** 日    期: 2022年2月13日18:01:35
**********************************************************/
func GetUrlDifferentRes(videoName, localFileName string, vid int, oss bool) (map[string]string, int) {
	urls := map[string]string{
		"res360":   "",
		"res480":   "",
		"res720":   "",
		"res1080":  "",
		"original": "",
	}
	ossDir := "video"
	if !oss {
		ossDir = "output"
	}
	maxRes, err := PreTreatmentVideo(localFileName, vid)
	maxRes = util.Min(maxRes, viper.GetInt("transcoding.max_res"))
	if err != nil {
		//调用审核失败
		VideoReviewFail(vid, "视频处理出现错误")
		return nil, 0
	}
	switch maxRes {
	case 1080:
		urls["res1080"] = GetUrl() + ossDir + "/" + videoName + "/1080p/" + "index.m3u8"
		fallthrough
	case 720:
		urls["res720"] = GetUrl() + ossDir + "/" + videoName + "/720p/" + "index.m3u8"
		fallthrough
	case 480:
		urls["res480"] = GetUrl() + ossDir + "/" + videoName + "/480p/" + "index.m3u8"
		fallthrough
	case 360:
		urls["res360"] = GetUrl() + ossDir + "/" + videoName + "/360p/" + "index.m3u8"
	}
	return urls, maxRes
}

/*********************************************************
** 函数功能: 创建不同分辨率文件夹
** 日    期: 2022年2月13日18:47:04
**********************************************************/
func CreateResDir(maxRes int, dirName string) {
	switch maxRes {
	case 1080:
		os.Mkdir("./file/output/"+dirName+"/1080p", os.ModePerm)
		fallthrough
	case 720:
		os.Mkdir("./file/output/"+dirName+"/720p", os.ModePerm)
		fallthrough
	case 480:
		os.Mkdir("./file/output/"+dirName+"/480p", os.ModePerm)
		fallthrough
	case 360:
		os.Mkdir("./file/output/"+dirName+"/360p", os.ModePerm)
	}
}

/*********************************************************
** 函数功能: 删除临时文件
** 日    期: 2022年2月13日19:00:44
**********************************************************/
func DeleteTempFile(maxRes int, dirName string) {
	if maxRes == 0 {
		os.Remove("./file/output/" + dirName + "/temp.m3u8")
		os.Remove("./file/output/" + dirName + "/temp_original.ts")
	} else {
		switch maxRes {
		case 1080:
			os.Remove("./file/output/" + dirName + "/temp_1080p.ts")
			os.Remove("./file/output/" + dirName + "/temp_1080p.mp4")
			fallthrough
		case 720:
			os.Remove("./file/output/" + dirName + "/temp_720p.ts")
			os.Remove("./file/output/" + dirName + "/temp_720p.mp4")
			fallthrough
		case 480:
			os.Remove("./file/output/" + dirName + "/temp_480p.ts")
			os.Remove("./file/output/" + dirName + "/temp_480p.mp4")
			fallthrough
		case 360:
			os.Remove("./file/output/" + dirName + "/temp_360p.ts")
			os.Remove("./file/output/" + dirName + "/temp_360p.mp4")
		}
	}
}
