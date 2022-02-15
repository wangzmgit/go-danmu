package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/common"
	"kuukaa.fun/danmu-v4/dto"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/service"
	"kuukaa.fun/danmu-v4/util"
)

/*********************************************************
** 函数功能: 上传视频信息
** 日    期: 2021/7/16
** 修改时间: 2021/10/31
** 版    本: 3.3.0
** 修改内容: 可以上传子视频信息
**********************************************************/
func UploadVideoInfo(ctx *gin.Context) {
	var video dto.UploadVideoDto
	err := ctx.Bind(&video)
	if err != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}
	title := video.Title
	cover := video.Cover
	partition := video.Partition
	uid, _ := ctx.Get("id")

	//验证数据
	if len(title) == 0 {
		response.CheckFail(ctx, nil, response.TitleCheck)
		return
	}
	if len(cover) == 0 {
		response.CheckFail(ctx, nil, response.CoverCheck)
		return
	}
	if partition == 0 {
		response.CheckFail(ctx, nil, response.PartitionCheck)
		return
	}

	res := service.UploadVideoInfoService(video, uid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 获取视频状态
** 日    期:2021/7/16
**********************************************************/
func GetVideoStatus(ctx *gin.Context) {
	vid, _ := strconv.Atoi(ctx.Query("vid"))
	uid, _ := ctx.Get("id")
	res := service.GetVideoStatusService(vid, uid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 修改视频信息
** 日    期:2021/7/17
**********************************************************/
func ModifyVideoInfo(ctx *gin.Context) {
	//获取参数
	var video dto.ModifyVideoDto
	err := ctx.Bind(&video)
	if err != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}
	title := video.Title
	cover := video.Cover

	if len(title) == 0 {
		response.CheckFail(ctx, nil, response.TitleCheck)
		return
	}
	if len(cover) == 0 {
		response.CheckFail(ctx, nil, response.CoverCheck)
		return
	}

	//从上下文中获取用户id
	uid, _ := ctx.Get("id")
	res := service.ModifyVideoInfoService(video, uid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 删除视频
** 日    期:2021/7/17
**********************************************************/
func DeleteVideo(ctx *gin.Context) {
	//获取参数
	var video dto.VideoIdDto
	err := ctx.Bind(&video)
	if err != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}
	id := video.ID
	uid, _ := ctx.Get("id")
	//数据验证
	if id == 0 {
		response.CheckFail(ctx, nil, response.VideoNotExist)
		return
	}
	//删除视频
	res := service.DeleteVideoService(id, uid)

	//删除播放量数据
	Redis := common.RedisClient
	if Redis != nil {
		Redis.Del(util.VideoClicksKey(int(id)))
	}

	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 获取自己的视频
** 日    期:2021/7/17
**********************************************************/
func GetMyUploadVideo(ctx *gin.Context) {
	uid, _ := ctx.Get("id")
	//获取分页信息
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	if page <= 0 || pageSize <= 0 {
		response.CheckFail(ctx, nil, response.PageOrSizeError)
		return
	}
	res := service.GetMyUploadVideoService(page, pageSize, uid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 通过ID获取视频
** 日    期: 2021/7/19
** 修改时间: 2021/10/31
** 版    本: 3.3.0
** 修改内容: 获取子视频列表
** 版    本: 3.5.0
** 修改内容: 移除子视频
**********************************************************/
func GetVideoByID(ctx *gin.Context) {
	vid, _ := strconv.Atoi(ctx.Query("vid"))
	if vid == 0 {
		response.CheckFail(ctx, nil, response.VideoNotExist)
		return
	}

	uid, _ := ctx.Get("uid")
	res := service.GetVideoByIDService(vid, ctx.ClientIP(), uid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 获取收藏列表
** 日    期:2021/7/22
**********************************************************/
func GetCollectVideo(ctx *gin.Context) {
	uid, _ := ctx.Get("id")
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))

	if page <= 0 || pageSize <= 0 {
		response.Fail(ctx, nil, response.PageOrSizeError)
		return
	}

	res := service.GetCollectVideoService(uid, page, pageSize)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 获取推荐视频
** 日    期:2021/8/1
** 修改时间: 2021/10/26
** 版    本: 3.3.0
** 修改内容: 获取合集所属的视频，不获取合集子视频
**********************************************************/
func GetRecommendVideo(ctx *gin.Context) {
	//因为视频比较少，就直接按播放量排名
	res := service.GetRecommendVideoService()
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 获取视频列表
** 日    期: 2021/8/1
** 修改时间: 2021/10/26
** 版    本: 3.3.0
** 修改内容: 获取合集所属的视频，不获取合集子视频
** 修改时间: 2021年12月11日17:03:06
** 版    本: 3.6.8
** 修改内容: 按分区获取视频列表
**********************************************************/
func GetVideoList(ctx *gin.Context) {
	var request dto.GetVideoListDto
	request.Page, _ = strconv.Atoi(ctx.Query("page"))
	request.PageSize, _ = strconv.Atoi(ctx.Query("page_size"))
	request.Partition, _ = strconv.Atoi(ctx.DefaultQuery("partition", "0")) //分区

	if request.Page <= 0 || request.PageSize <= 0 {
		response.Fail(ctx, nil, response.PageOrSizeError)
		return
	}
	if request.PageSize >= 30 {
		response.Fail(ctx, nil, response.TooManyRequests)
		return
	}

	res := service.GetVideoListService(request)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 通过用户ID获取视频列表
** 日    期:2021/8/4
** 修改时间: 2021/10/26
** 版    本: 3.3.0
** 修改内容: 获取合集所属的视频，不获取合集子视频
**********************************************************/
func GetVideoListByUserID(ctx *gin.Context) {
	uid, _ := strconv.Atoi(ctx.Query("uid"))
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))

	if page <= 0 || pageSize <= 0 {
		response.Fail(ctx, nil, response.PageOrSizeError)
		return
	}

	res := service.GetVideoListByUserIDService(uid, page, pageSize)
	response.HandleResponse(ctx, res)
}
