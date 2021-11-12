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
** 函数功能: 获取审核列表
** 日    期: 2021年11月12日14:52:09
**********************************************************/
func GetReviewVideoListService(page int, pageSize int) response.ResponseStruct {
	var total int //记录总数
	var videos []model.Video

	DB := common.GetDB()
	DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
	DB.Raw("select * from videos where deleted_at is null and id in (select vid from reviews where deleted_at is null and status = 1000)").Scan(&videos)

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"count": total, "videos": vo.ToAdminVideoVo(videos)},
		Msg:        "ok",
	}
}

/*********************************************************
** 函数功能: 审核视频
** 日    期:2021/8/4
**********************************************************/
func ReviewVideoService(review dto.ReviewRequest, isReview bool) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}

	DB := common.GetDB()
	tx := DB.Begin()
	if err := tx.Model(&model.Video{}).Where("id = ?", review.VID).Updates(
		map[string]interface{}{
			"review": isReview,
		},
	).Error; err != nil {
		tx.Rollback()
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = "更新视频状态失败"
		return res
	}
	//创建审核状态
	if err := tx.Model(&model.Review{}).Where("vid = ?", review.VID).Updates(
		map[string]interface{}{
			"status":  review.Status,
			"remarks": review.Remarks,
		},
	).Error; err != nil {
		tx.Rollback()
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = "更新审核状态失败"
		return res
	}
	tx.Commit()
	return res
}
