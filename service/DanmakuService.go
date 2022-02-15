package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/common"
	"kuukaa.fun/danmu-v4/dto"
	"kuukaa.fun/danmu-v4/model"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/vo"
)

func GetDanmakuService(vid int) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}

	var danmakuList []vo.DanmakuVo
	DB := common.GetDB()
	if !IsVideoExist(DB, uint(vid)) {
		res.HttpStatus = http.StatusUnprocessableEntity
		res.Code = response.CheckFailCode
		res.Msg = response.VideoNotExist
		return res
	}
	DB.Model(&model.Danmaku{}).Select("time,type,color,text").Where("vid = ? ", vid).Scan(&danmakuList)
	res.Data = gin.H{"danmaku": danmakuList}
	return res
}

func SendDanmaku(danmaku dto.DanmakuDto, uid interface{}) response.ResponseStruct {
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
		Msg:        response.OK,
	}
}
