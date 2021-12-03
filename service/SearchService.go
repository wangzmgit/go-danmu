package service

import (
	"net/http"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/vo"

	"github.com/gin-gonic/gin"
)

/*********************************************************
** 函数功能: 搜索
** 日    期:2021/11/11
**********************************************************/
func SearchService(keywords string) response.ResponseStruct {
	var videos []vo.SearchVideoVo
	DB := common.GetDB()
	DB = DB.Limit(50)
	DB.Model(model.Video{}).Select("id,title,cover").Where("title like ?", keywords).Scan(&videos)
	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"videos": videos},
		Msg:        "ok",
	}
}
