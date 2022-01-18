package manage

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/dto"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/service"
)

/*********************************************************
** 函数功能: 获取视频列表
** 日    期: 2021/8/4
**********************************************************/
func AdminGetVideoList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	videoFrom := ctx.DefaultQuery("video_from", "user")
	if page <= 0 || pageSize <= 0 {
		response.CheckFail(ctx, nil, "页码或数量有误")
		return
	}

	res := service.AdminGetVideoListService(page, pageSize, videoFrom)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 删除视频
** 日    期:2021/8/3
**********************************************************/
func AdminDeleteVideo(ctx *gin.Context) {
	var request dto.AdminIdDto
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	id := request.ID

	if id == 0 {
		response.CheckFail(ctx, nil, "视频不存在")
		return
	}

	res := service.AdminDeleteVideoService(id)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 管理员导入视频
** 日    期: 2021/10/6
**********************************************************/
func ImportVideo(ctx *gin.Context) {
	var video dto.ImportVideo
	err := ctx.Bind(&video)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	title := video.Title
	cover := video.Cover

	//验证数据
	if video.Type != "mp4" && video.Type != "hls" {
		response.CheckFail(ctx, nil, "视频类型错误")
		return
	}
	if len(title) == 0 {
		response.CheckFail(ctx, nil, "标题不能为空")
		return
	}
	if len(cover) == 0 {
		response.CheckFail(ctx, nil, "封面图不能为空")
		return
	}

	res := service.ImportVideoService(video)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 管理员导入视频分集资源
** 日    期: 2022年1月13日16:21:48
**********************************************************/
func ImportResource(ctx *gin.Context) {
	var video dto.ImportResourceDto
	err := ctx.Bind(&video)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}

	if video.Vid == 0 {
		response.CheckFail(ctx, nil, "视频不存在")
		return
	}
	if len(video.Original) == 0 && len(video.Res360) == 0 {
		response.CheckFail(ctx, nil, "视频链接不存在")
		return
	}

	res := service.ImportResourceService(video)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 管理员获取视频资源
** 日    期: 2022年1月14日11:27:32
**********************************************************/
func GetResourceList(ctx *gin.Context) {
	vid, _ := strconv.Atoi(ctx.Query("vid"))
	if vid <= 0 {
		response.CheckFail(ctx, nil, "视频不存在")
		return
	}

	res := service.GetResourceListService(vid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 管理员删除视频资源
** 日    期: 2022年1月14日11:27:32
**********************************************************/
func DeleteResource(ctx *gin.Context) {
	var id dto.UUID
	err := ctx.Bind(&id)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}

	res := service.DeleteResourceService(id.UUID)
	response.HandleResponse(ctx, res)
}
