package controller

import (
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/service"

	"github.com/gin-gonic/gin"
)

/*********************************************************
** 函数功能: 创建合集
** 日    期: 2021年11月19日14:48:29
** 版    本: 3.6.0
**********************************************************/
func CreateCollection(ctx *gin.Context) {
	var request dto.CreateCollectionDto
	err := ctx.Bind(&request)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	title := request.Title
	cover := request.Cover
	uid, _ := ctx.Get("id")

	//验证数据
	if len(title) == 0 {
		response.CheckFail(ctx, nil, "标题不能为空")
		return
	}

	if len(cover) == 0 {
		response.CheckFail(ctx, nil, "封面图不能为空")
		return
	}

	res := service.CreateCollectionService(request, uid)
	response.HandleResponse(ctx, res)
}
