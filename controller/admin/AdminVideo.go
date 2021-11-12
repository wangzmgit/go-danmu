package admin_controller

import (
	"strconv"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/service"

	"github.com/gin-gonic/gin"
)

/*********************************************************
** 函数功能: 获取视频列表
** 日    期: 2021/8/4
**********************************************************/
func AdminGetVideoList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	if page <= 0 || pageSize <= 0 {
		response.CheckFail(ctx, nil, "页码或数量有误")
		return
	}

	res := service.AdminGetVideoListService(page, pageSize)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 删除视频
** 日    期:2021/8/3
**********************************************************/
func AdminDeleteVideo(ctx *gin.Context) {
	var request dto.AdminIDRequest
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
	var request dto.ImportVideo
	err := ctx.Bind(&request)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	title := request.Title
	cover := request.Cover
	video := request.Video

	//验证数据
	if len(title) == 0 {
		response.CheckFail(ctx, nil, "标题不能为空")
		return
	}
	if len(cover) == 0 {
		response.CheckFail(ctx, nil, "封面图不能为空")
		return
	}
	if len(video) == 0 {
		response.CheckFail(ctx, nil, "视频链接不能为空")
		return
	}

	res := service.ImportVideoService(request)
	response.HandleResponse(ctx, res)
}
