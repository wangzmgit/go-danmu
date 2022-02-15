package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/dto"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/service"
)

func GetDanmaku(ctx *gin.Context) {
	vid, _ := strconv.Atoi(ctx.Query("vid"))
	if vid == 0 {
		response.Fail(ctx, nil, response.ParameterError)
		return
	}

	res := service.GetDanmakuService(vid)
	response.HandleResponse(ctx, res)
}

func SendDanmaku(ctx *gin.Context) {
	var danmaku dto.DanmakuDto
	err := ctx.ShouldBind(&danmaku)
	if err != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}
	//内容
	vid := danmaku.Vid
	time := danmaku.Time
	text := danmaku.Text
	uid, _ := ctx.Get("id")
	if vid == 0 || time == 0 {
		response.CheckFail(ctx, nil, response.SendFail)
		return
	}
	if text == "" {
		response.CheckFail(ctx, nil, response.DanmakuCheck)
		return
	}

	res := service.SendDanmaku(danmaku, uid)
	response.HandleResponse(ctx, res)
}
