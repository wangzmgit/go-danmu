package service

import (
	"net/http"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/vo"

	"github.com/gin-gonic/gin"
)

func GetDanmakuService(vid int) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}

	var danmakuList []vo.DanmakuVo
	DB := common.GetDB()
	if !IsVideoExist(DB, uint(vid)) {
		res.HttpStatus = http.StatusUnprocessableEntity
		res.Code = response.CheckFailCode
		res.Msg = "找不到视频"
		return res
	}
	DB.Model(&model.Danmaku{}).Select("time,type,color,text").Where("vid = ? ", vid).Scan(&danmakuList)
	res.Data = gin.H{"danmaku": danmakuList}
	return res
}

func SendDanmaku(danmaku dto.DanmakuRequest, uid interface{}) response.ResponseStruct {
	DB := common.GetDB()
	newDanmaku := model.Danmaku{
		Vid:   danmaku.Vid,
		Time:  danmaku.Time,
		Type:  danmaku.Type,
		Color: danmaku.Color,
		Text:  danmaku.Text,
		Uid:   uid.(uint),
	}
	DB.Create(&newDanmaku)
	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}
}
