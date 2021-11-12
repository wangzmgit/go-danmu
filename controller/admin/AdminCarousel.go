package admin_controller

import (
	"os"
	"path"
	"strconv"
	"time"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/service"
	"wzm/danmu3.0/util"

	"github.com/gin-gonic/gin"
)

/*********************************************************
** 函数功能: 上传轮播图
** 日    期: 2021/8/4
**********************************************************/
func UploadCarousel(ctx *gin.Context) {
	carousel, err := ctx.FormFile("carousel")
	if err != nil {
		response.Fail(ctx, nil, "图片上传失败")
		return
	}
	suffix := path.Ext(carousel.Filename)
	if suffix != ".jpg" && suffix != ".jpeg" && suffix != ".png" {
		response.CheckFail(ctx, nil, "图片不符合要求")
		return
	}
	carousel.Filename = util.RandomString(3) + strconv.FormatInt(time.Now().UnixNano(), 10) + suffix
	errSave := ctx.SaveUploadedFile(carousel, "./file/carousel/"+carousel.Filename)
	if errSave != nil {
		response.Fail(ctx, nil, "图片保存失败")
		return
	}
	fileInfo, err := os.Stat("./file/carousel/" + carousel.Filename)
	//大小限制到5M
	if fileInfo == nil || fileInfo.Size() > 1024*1024*5 || err != nil {
		response.CheckFail(ctx, nil, "图片大小不符合要求")
		return
	}
	// 拼接上传图片的路径信息
	localFileName := "./file/carousel/" + carousel.Filename
	objectName := "carousel/" + carousel.Filename

	res := service.UploadCarouselService(localFileName, objectName)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 上传轮播图信息
** 日    期: 2021/8/4
**********************************************************/
func UploadCarouselInfo(ctx *gin.Context) {
	//获取参数
	var carousel dto.CarouselRequest
	err := ctx.Bind(&carousel)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	img := carousel.Img
	url := carousel.Url

	//验证数据
	if len(img) == 0 {
		response.CheckFail(ctx, nil, "图片不能为空")
		return
	}

	res := service.UploadCarouselInfoService(img, url)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 获取轮播图
** 日    期: 2021/8/4
**********************************************************/
func AdminGetCarousel(ctx *gin.Context) {
	res := service.AdminGetCarouselService()
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 删除轮播图
** 日    期: 2021/8/4
**********************************************************/
func DeleteCarousel(ctx *gin.Context) {
	//获取参数
	var request dto.AdminIDRequest
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	id := request.ID

	if id == 0 {
		response.CheckFail(ctx, nil, "轮播图不存在")
		return
	}

	res := service.DeleteCarouselService(id)
	response.HandleResponse(ctx, res)
}
