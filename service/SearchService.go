package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/common"
	"kuukaa.fun/danmu-v4/model"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/vo"
)

/*********************************************************
** 函数功能: 搜索
** 日    期:2021/11/11
**********************************************************/
func SearchService(keywords string) response.ResponseStruct {
	var videos []vo.SearchVideoVo
	DB := common.GetDB()
	DB = DB.Limit(50)
	DB.Model(model.Video{}).Select("id,title,cover").Where("title like ? and review = 1", keywords).Scan(&videos)
	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"videos": videos},
		Msg:        response.OK,
	}
}
