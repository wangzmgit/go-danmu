package admin_controller

import (
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/service"

	"github.com/gin-gonic/gin"
)

/*********************************************************
** 函数功能: 获取公告
** 日    期: 2021/8/4
**********************************************************/
func AdminGetAnnounce(ctx *gin.Context) {
	res := service.AdminGetAnnounceService()
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 添加公告
** 日    期: 2021/8/4
**********************************************************/
func AddAnnounce(ctx *gin.Context) {
	//获取参数
	var announce dto.AddAnnounceRequest
	err := ctx.Bind(&announce)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	title := announce.Title
	content := announce.Content

	if len(title) == 0 {
		response.CheckFail(ctx, nil, "标题不能为空")
		return
	}
	if len(content) == 0 {
		response.CheckFail(ctx, nil, "内容不能为空")
		return
	}

	res := service.AddAnnounceService(announce)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 删除公告
** 日    期: 2021/8/4
**********************************************************/
func DeleteAnnounce(ctx *gin.Context) {
	var request dto.AdminIDRequest
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	id := request.ID

	res := service.DeleteAnnounceService(id)
	response.HandleResponse(ctx, res)
}
