package service

import (
	"net/http"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/vo"

	"github.com/gin-gonic/gin"
)

/*********************************************************
** 函数功能: 创建反馈(站内)
** 日    期: 2021/12/3
**********************************************************/
func CreateOpinionOnSiteService(desc string, uid uint) response.ResponseStruct {
	DB := common.GetDB()

	newOpinion := model.Opinion{
		Desc: desc,
		Uid:  uid,
	}

	DB.Create(&newOpinion)
	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}
}

/*********************************************************
** 函数功能: 创建反馈
** 日    期: 2021/12/3
**********************************************************/
func CreateOpinionService(opinion dto.OpinionDto) response.ResponseStruct {
	DB := common.GetDB()

	newOpinion := model.Opinion{
		Name:      opinion.Name,
		Email:     opinion.Email,
		Telephone: opinion.Telephone,
		Gender:    opinion.Gender,
		Desc:      opinion.Desc,
	}

	DB.Create(&newOpinion)
	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}
}

/*********************************************************
** 函数功能: 获取反馈列表
** 日    期: 2021/12/3
**********************************************************/
func GetOpinionListService(page int, pageSize int) response.ResponseStruct {
	var total int //记录总数
	var opinions []vo.OpinionVo
	DB := common.GetDB()
	DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
	DB.Model(&model.Opinion{}).Select("id,name,email,telephone,gender,`desc`,uid,created_at").Scan(&opinions).Count(&total)

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"count": total, "opinions": opinions},
		Msg:        "ok",
	}
}
