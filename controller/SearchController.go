package controller

import (
	"strings"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/service"
	"wzm/danmu3.0/util"

	"github.com/gin-gonic/gin"
)

/*********************************************************
** 函数功能: 搜索
** 日    期:
**********************************************************/
func Search(ctx *gin.Context) {
	search := ctx.DefaultQuery("keywords", "0")
	if search == "0" || util.ExistSQLInject(search) {
		response.CheckFail(ctx, nil, "请输入搜索内容")
		return
	}

	//拆分关键词
	keywordsList := strings.Fields(search)
	length := len(keywordsList)
	if length > 5 {
		response.CheckFail(ctx, nil, "输入的关键词过多")
		return
	}
	//拼接查询语句
	keywords := "title like '%" + keywordsList[0] + "%'"
	for i := 1; i < length; i++ {
		keywords += "or title like '%" + keywordsList[i] + "%'"
	}

	res := service.SearchService(keywords)
	response.HandleResponse(ctx, res)
}
