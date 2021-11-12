package controller

import (
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/service"

	"github.com/gin-gonic/gin"
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
