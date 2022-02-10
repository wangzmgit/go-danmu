package manage

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/dto"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/service"
)

/*********************************************************
** 函数功能: 获取待审核视频列表
** 日    期: 2021年11月12日14:55:29
**********************************************************/
func GetReviewVideoList(ctx *gin.Context) {
	//获取参数
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))

	if page <= 0 || pageSize <= 0 {
		response.CheckFail(ctx, nil, "页码或数量有误")
		return
	}

	res := service.GetReviewVideoListService(page, pageSize)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 审核视频
** 日    期: 2021年11月12日15:01:04
**********************************************************/
func ReviewVideo(ctx *gin.Context) {
	//获取参数
	var reviewRequest dto.ReviewDto
	err := ctx.Bind(&reviewRequest)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	vid := reviewRequest.VID
	status := reviewRequest.Status
	var isReview bool
	if vid == 0 {
		response.CheckFail(ctx, nil, "视频不存在")
		return
	}

	if status == 2000 {
		isReview = true
	} else if status == 4001 || status == 4002 {
		isReview = false
	} else {
		response.CheckFail(ctx, nil, "状态错误")
		return
	}

	res := service.ReviewVideoService(reviewRequest, isReview)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 通过视频ID获取待审核视频资源
** 日    期: 2022年1月10日11:02:51
**********************************************************/
func GetReviewVideoByID(ctx *gin.Context) {
	vid, _ := strconv.Atoi(ctx.Query("vid"))
	if vid == 0 {
		response.CheckFail(ctx, nil, "视频不见了")
		return
	}

	res := service.GetReviewVideoByIDService(vid)
	response.HandleResponse(ctx, res)
}