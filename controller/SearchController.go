package controller

import (
	"strings"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/util"

	"github.com/gin-gonic/gin"
)

func Search(ctx *gin.Context) {
	search := ctx.DefaultQuery("keywords", "0")
	if search == "0" || util.ExistSQLInject(search) {
		response.CheckFail(ctx, nil, "请输入搜索内容")
		return
	}
	var videos []dto.SearchVideoDto
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

	DB := common.GetDB()
	DB = DB.Limit(50)
	sql := "select id,title,cover from videos where review = true and " + keywords
	DB.Raw(sql).Scan(&videos)
	response.Success(ctx, gin.H{"videos": videos}, "ok")

}
