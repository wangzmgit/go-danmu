package controller

import (
	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/service"
)

/*********************************************************
** 函数功能: 获取公告
** 日    期:2021/7/29
**********************************************************/
func GetAnnounce(ctx *gin.Context) {
	uid, _ := ctx.Get("id")
	res := service.GetAnnounceService(uid)
	response.HandleResponse(ctx, res)
}
