package controller

import (
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/service"

	"github.com/gin-gonic/gin"
)

/*********************************************************
** 函数功能: 获取轮播图
** 日    期:
**********************************************************/
func GetCarousel(ctx *gin.Context) {
	res := service.GetCarouselService()
	response.HandleResponse(ctx, res)
}
