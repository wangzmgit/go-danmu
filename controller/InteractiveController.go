package controller

import (
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/service"

	"github.com/gin-gonic/gin"
)

/*********************************************************
** 函数功能: 添加收藏
** 日    期:2021/7/22
**********************************************************/
func Collect(ctx *gin.Context) {
	//获取参数
	var request dto.InteractiveRequest
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}

	vid := request.ID
	uid, _ := ctx.Get("id")
	//验证数据
	if vid <= 0 {
		response.CheckFail(ctx, nil, "视频不见了")
		return
	}

	res := service.CollectService(vid, uid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 取消收藏
** 日    期:2021/7/22
**********************************************************/
func CancelCollect(ctx *gin.Context) {
	//获取参数
	var request dto.InteractiveRequest
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}

	vid := request.ID
	uid, _ := ctx.Get("id")

	res := service.CancelCollectService(vid, uid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 点赞
** 日    期:2021/7/22
**********************************************************/
func Like(ctx *gin.Context) {
	//获取参数
	var request dto.InteractiveRequest
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	vid := request.ID
	uid, _ := ctx.Get("id")
	//验证数据
	if vid <= 0 {
		response.CheckFail(ctx, nil, "视频不见了")
		return
	}

	res := service.LikeService(vid, uid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 取消赞
** 日    期:2021/7/22
**********************************************************/
func Dislike(ctx *gin.Context) {
	//获取参数
	var request dto.InteractiveRequest
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	vid := request.ID
	uid, _ := ctx.Get("id")
	res := service.DislikeService(vid, uid)
	response.HandleResponse(ctx, res)
}
