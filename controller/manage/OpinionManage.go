package manage

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/service"
)

/*********************************************************
** 函数功能: 获取意见反馈列表
** 日    期: 2021年11月12日14:55:29
**********************************************************/
func GetOpinionList(ctx *gin.Context) {
	//获取参数
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))

	if page <= 0 || pageSize <= 0 {
		response.CheckFail(ctx, nil, "页码或数量有误")
		return
	}

	if pageSize >= 30 {
		response.CheckFail(ctx, nil, "请求数量过多")
		return
	}

	res := service.GetOpinionListService(page, pageSize)
	response.HandleResponse(ctx, res)
}
