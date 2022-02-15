package controller

import (
	"strings"

	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/service"
)

/*********************************************************
** 函数功能: 搜索
** 日    期:
** 修改时间: 2021年11月27日10:55:32
** 版    本: 3.6.4
** 修改内容: 修改关键词处理
**********************************************************/
func Search(ctx *gin.Context) {
	search := ctx.Query("keywords")
	if len(search) == 0 {
		response.CheckFail(ctx, nil, response.SearchCheck)
		return
	}

	keywords := "%" + strings.Replace(search, " ", "%", -1) + "%"

	res := service.SearchService(keywords)
	response.HandleResponse(ctx, res)
}
