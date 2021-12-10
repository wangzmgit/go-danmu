package controller

import (
	"strconv"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/service"

	"github.com/gin-gonic/gin"
)

/*********************************************************
** 函数功能: 获取分区
** 日    期: 2021年12月9日13:39:33
**********************************************************/
func GetPartitionList(ctx *gin.Context) {
	fid, _ := strconv.Atoi(ctx.DefaultQuery("fid", "0"))
	if fid < 0 {
		response.Fail(ctx, nil, "参数有误")
		return
	}

	res := service.GetPartitionListService(fid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 获取所有分区
** 日    期: 2021年12月9日20:50:56
**********************************************************/
func GetAllPartition(ctx *gin.Context) {
	res := service.GetAllPartitionService()
	response.HandleResponse(ctx, res)
}
