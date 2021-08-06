package controller

import (
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"

	"github.com/gin-gonic/gin"
)

func GetCarousel(ctx *gin.Context) {
	DB := common.GetDB()
	var carousels []dto.CarouselDto
	DB.Model(&model.Carousel{}).Select("img,url").Scan(&carousels)
	response.Success(ctx, gin.H{"carousels": carousels}, "ok")
}
