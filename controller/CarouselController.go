package controller

import (
	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/service"
)

/*********************************************************
** 函数功能: 获取轮播图
** 日    期:
**********************************************************/
func GetCarousel(ctx *gin.Context) {
	res := service.GetCarouselService()
	response.HandleResponse(ctx, res)
}
