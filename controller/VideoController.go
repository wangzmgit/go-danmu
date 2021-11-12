package controller

import (
	"strconv"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/service"
	"wzm/danmu3.0/util"

	"github.com/gin-gonic/gin"
)

/*********************************************************
** 函数功能: 上传视频信息
** 日    期:2021/7/16
** 修改时间: 2021/10/31
** 版    本: 3.3.0
** 修改内容: 可以上传子视频信息
**********************************************************/
func UploadVideoInfo(ctx *gin.Context) {
	var video dto.UploadVideoRequest
	err := ctx.Bind(&video)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	title := video.Title
	cover := video.Cover
	parent := video.Parent
	uid, _ := ctx.Get("id")

	//验证数据
	if len(title) == 0 {
		response.CheckFail(ctx, nil, "标题不能为空")
		return
	}
	if parent == 0 && len(cover) == 0 {
		response.CheckFail(ctx, nil, "封面图不能为空")
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
	var video dto.VideoModifyRequest
	err := ctx.Bind(&video)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	title := video.Title
	cover := video.Cover

	if len(title) == 0 {
		response.CheckFail(ctx, nil, "标题不能为空")
		return
	}
	if len(cover) == 0 {
		response.CheckFail(ctx, nil, "封面图不能为空")
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
	var video dto.DeleteVideoRequest
	err := ctx.Bind(&video)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	id := video.ID
	uid, _ := ctx.Get("id")
	//数据验证
	if id == 0 {
		response.CheckFail(ctx, nil, "视频不存在")
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
		response.CheckFail(ctx, nil, "页码或数量有误")
		return
	}
	res := service.GetMyUploadVideoService(page, pageSize, uid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 视频信息修改请求
** 日    期:2021/7/18
**********************************************************/
func UpdateRequest(ctx *gin.Context) {
	var review dto.UpdateVideoReviewRequest
	err := ctx.Bind(&review)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	status := review.Status
	uid, _ := ctx.Get("id")
	if status == 5001 || status == 5002 {
		res := service.UpdateRequestService(review, uid)
		response.HandleResponse(ctx, res)
	} else {
		response.Fail(ctx, nil, "申请状态有误")
	}
}

/*********************************************************
** 函数功能: 通过ID获取视频
** 日    期: 2021/7/19
** 修改时间: 2021/10/31
** 版    本: 3.3.0
** 修改内容: 获取子视频列表
**********************************************************/
func GetVideoByID(ctx *gin.Context) {
	vid, _ := strconv.Atoi(ctx.Query("vid"))
	if vid == 0 {
		response.CheckFail(ctx, nil, "视频不见了")
		return
	}
	res := service.GetVideoByIDService(vid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 获取视频交互数据
** 日    期:2021/7/22
**********************************************************/
func GetVideoInteractiveData(ctx *gin.Context) {
	uid, _ := ctx.Get("id")
	vid, _ := strconv.Atoi(ctx.Query("vid"))

	if vid == 0 {
		response.CheckFail(ctx, nil, "视频不见了")
		return
	}

	res := service.GetVideoInteractiveDataService(vid, uid)
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
		response.Fail(ctx, nil, "页码或数量有误")
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
** 日    期:2021/8/1
** 修改时间: 2021/10/26
** 版    本: 3.3.0
** 修改内容: 获取合集所属的视频，不获取合集子视频
**********************************************************/
func GetVideoList(ctx *gin.Context) {

	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))

	if page <= 0 || pageSize <= 0 {
		response.Fail(ctx, nil, "页码或数量有误")
		return
	}

	res := service.GetVideoListService(page, pageSize)
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
		response.Fail(ctx, nil, "页码或数量有误")
		return
	}

	res := service.GetVideoListByUserIDService(uid, page, pageSize)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 通过视频ID获取子视频列表
** 日    期:2021/11/6
**********************************************************/
func GetSubVideoListByVideoID(ctx *gin.Context) {
	//获取参数
	uid, _ := ctx.Get("id")
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	parentId, _ := strconv.Atoi(ctx.Query("parent_id"))
	if page <= 0 || pageSize <= 0 || parentId <= 0 {
		response.Fail(ctx, nil, "请求参数有误")
		return
	}

	res := service.GetSubVideoListByVideoIDService(uid, page, pageSize, parentId)
	response.HandleResponse(ctx, res)
}
