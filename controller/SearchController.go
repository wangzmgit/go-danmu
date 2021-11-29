package controller

import (
	"strings"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/service"

	"github.com/gin-gonic/gin"
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
		response.CheckFail(ctx, nil, "请输入搜索内容")
		return
	}

	keywords := "%" + strings.Replace(search, " ", "%", -1) + "%"

	res := service.SearchService(keywords)
	response.HandleResponse(ctx, res)
}
