package admin_controller

import (
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/service"

	"github.com/gin-gonic/gin"
)

/*********************************************************
** 函数功能: 添加分区
** 日    期: 2021年12月9日17:27:07
**********************************************************/
func AddPartition(ctx *gin.Context) {
	//获取参数
	var partition dto.PartitionDto
	err := ctx.Bind(&partition)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}

	content := partition.Content

	if len(content) == 0 {
		response.CheckFail(ctx, nil, "内容不能为空")
		return
	}

	res := service.AddPartitionService(partition)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 删除分区
** 日    期: 2021年12月9日17:39:38
**********************************************************/
func DeletePartition(ctx *gin.Context) {
	var request dto.DeletePartitionDto
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	id := request.ID

	if id == 0 {
		response.CheckFail(ctx, nil, "分区不存在")
		return
	}

	res := service.DeletePartitionService(id)
	response.HandleResponse(ctx, res)
}
