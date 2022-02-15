package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/common"
	"kuukaa.fun/danmu-v4/dto"
	"kuukaa.fun/danmu-v4/model"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/vo"
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
	//统计数量
	DB.Model(&model.Review{}).Where("status = 1000").Count(&total)
	DB.Raw("select * from videos where deleted_at is null and id in (select vid from reviews where deleted_at is null and status = 1000)").Scan(&videos)

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"count": total, "videos": vo.ToAdminVideoListVo(videos)},
		Msg:        response.OK,
	}
}

/*********************************************************
** 函数功能: 通过视频ID获取待审核视频资源
** 日    期: 2022年1月6日15:16:50
**********************************************************/
func GetReviewVideoByIDService(vid int) response.ResponseStruct {
	//暂时不支持上传分P
	var videos model.Resource
	DB := common.GetDB()
	DB.Model(&model.Resource{}).Where("vid = ?", vid).Last(&videos)
	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"video": vo.ToReviewResourceVo(videos)},
		Msg:        response.OK,
	}
}

/*********************************************************
** 函数功能: 审核视频
** 日    期:2021/8/4
**********************************************************/
func ReviewVideoService(review dto.ReviewDto, isReview bool) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
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
		res.Msg = response.UpdateStatusFail
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
		res.Msg = response.UpdateStatusFail
		return res
	}
	tx.Commit()
	return res
}

/*********************************************************
** 函数功能: 视频处理失败未通过审核
** 日    期: 2022年1月5日17:12:27
**********************************************************/
func VideoReviewFail(vid int, msg string) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}

	DB := common.GetDB()
	tx := DB.Begin()
	if err := tx.Model(&model.Video{}).Where("id = ?", vid).Updates(
		map[string]interface{}{
			"review": false,
		},
	).Error; err != nil {
		tx.Rollback()
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = response.UpdateStatusFail
		return res
	}
	//创建审核状态
	if err := tx.Model(&model.Review{}).Where("vid = ?", vid).Updates(
		map[string]interface{}{
			"status":  4001,
			"remarks": msg,
		},
	).Error; err != nil {
		tx.Rollback()
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = response.UpdateStatusFail
		return res
	}
	tx.Commit()
	return res
}
