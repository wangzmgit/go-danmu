package controller

import (
	"net/http"
	"strconv"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"

	"github.com/gin-gonic/gin"
)

func GetDanmaku(ctx *gin.Context) {
	DB := common.GetDB()
	var danmakuList []dto.DanmakuDto
	vid, _ := strconv.Atoi(ctx.Query("vid"))
	if vid == 0 || !IsVideoExist(DB, uint(vid)) {
		response.Fail(ctx, nil, "视频不存在")
		return
	}
	DB.Model(&model.Danmaku{}).Select("time,type,color,text").Where("vid = ? ", vid).Scan(&danmakuList)
	response.Success(ctx, gin.H{"danmaku": danmakuList}, "ok")
}

func SendDanmaku(ctx *gin.Context) {
	DB := common.GetDB()
	var danmaku = model.Danmaku{}
	err := ctx.ShouldBind(&danmaku)
	if err != nil {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "请求错误")
		return
	}
	//内容
	vid := danmaku.Vid
	time := danmaku.Time
	danmakuType := danmaku.Type
	color := danmaku.Color
	text := danmaku.Text
	uid, _ := ctx.Get("id")
	if vid == 0 || time == 0 {
		response.CheckFail(ctx, nil, "发送失败")
		return
	}
	if text == "" {
		response.CheckFail(ctx, nil, "不能发送空内容")
		return
	}
	newDanmaku := model.Danmaku{
		Vid:   vid,
		Time:  time,
		Type:  danmakuType,
		Color: color,
		Text:  text,
		Uid:   uid.(uint),
	}
	DB.Create(&newDanmaku)
	response.Success(ctx, nil, "ok")
}
