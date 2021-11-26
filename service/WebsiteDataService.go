package service

import (
	"bytes"
	"net/http"
	"os/exec"
	"time"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/vo"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

/*********************************************************
** 函数功能: 获取网站数据
** 日    期: 2021年11月25日
**********************************************************/
func GetTotalWebsiteDataService() response.ResponseStruct {
	var data vo.TotalData
	DB := common.GetDB()

	DB.Model(&model.User{}).Count(&data.User)
	DB.Model(&model.Video{}).Count(&data.Video)
	DB.Model(&model.Review{}).Where("status = 1000").Count(&data.Review)

	//版本、redis、ffmpeg
	data.Version = common.Version
	data.Redis = GetRedisStatus()
	data.FFmpeg = GetFFmpegStatus()

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"data": data},
		Msg:        "ok",
	}
}

/*********************************************************
** 函数功能: 获取网站近期(5天)数据
** 日    期: 2021年11月25日
**********************************************************/
func GetGetRecentWebsiteDataService() response.ResponseStruct {
	DB := common.GetDB()
	data := make([]vo.OneDayData, 5)

	for i := 0; i < 5; i++ {
		data[i] = GetOneDayData(DB, i)
	}

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"data": data},
		Msg:        "ok",
	}
}

/*********************************************************
** 函数功能: 获取一天的视频和用户
** 日    期: 2021年11月25日
**********************************************************/
func GetOneDayData(db *gorm.DB, offset int) vo.OneDayData {
	var data vo.OneDayData

	t := time.Now()
	startTime := t.AddDate(0, 0, -(offset + 1))
	endTime := t.AddDate(0, 0, -offset)
	data.Date = startTime.Format("2006/1/02")

	db.Model(&model.User{}).Where("created_at > ? and created_at < ?", startTime, endTime).Count(&data.User)
	db.Model(&model.Video{}).Where("created_at > ? and created_at < ?", startTime, endTime).Count(&data.Video)

	return data
}

/*********************************************************
** 函数功能: 获取redis状态
** 日    期: 2021年11月26日
**********************************************************/
func GetRedisStatus() bool {
	if common.RedisClient == nil {
		return false
	}

	return true
}

/*********************************************************
** 函数功能: 获取ffmpeg状态
** 日    期: 2021年11月26日
**********************************************************/
func GetFFmpegStatus() bool {
	cmd := exec.Command("ffmpeg", "-version")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}
