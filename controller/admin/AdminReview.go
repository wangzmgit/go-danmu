package admin_controller

import (
	"strconv"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"

	"github.com/gin-gonic/gin"
)

/*********************************************************
** 函数功能: 获取待审核视频列表
** 日    期:2021/8/4
**********************************************************/
func GetReviewVideoList(ctx *gin.Context) {
	DB := common.GetDB()
	var videos []model.Video
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	if page > 0 && pageSize > 0 {
		//记录总数
		var total int
		DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
		DB.Raw("select * from videos where deleted_at is null and id in (select vid from reviews where deleted_at is null and status = 1000)").Scan(&videos)
		response.Success(ctx, gin.H{"count": total, "videos": dto.ToAdminVideoDto(videos)}, "ok")
	} else {
		response.Fail(ctx, nil, "获取数量有误")
	}
}

/*********************************************************
** 函数功能: 审核视频
** 日    期:2021/8/4
**********************************************************/
func ReviewVideo(ctx *gin.Context) {
	type review struct {
		VID     uint
		Status  int
		Remarks string
	}
	var requestReview review
	err := ctx.Bind(&requestReview)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	vid := requestReview.VID
	status := requestReview.Status
	remarks := requestReview.Remarks
	var isReview bool
	if vid == 0 {
		response.CheckFail(ctx, nil, "视频不存在")
		return
	}
	if status == 2000 {
		isReview = true
	} else if status == 4001 || status == 4002 {
		isReview = false
	} else {
		response.CheckFail(ctx, nil, "状态错误")
		return
	}
	DB := common.GetDB()
	tx := DB.Begin()
	if err := tx.Model(&model.Video{}).Where("id = ?", vid).Updates(map[string]interface{}{"review": isReview}).Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "修改失败")
		return
	}
	//创建审核状态
	if err := tx.Model(&model.Review{}).Where("vid = ?", vid).Updates(map[string]interface{}{"status": status, "remarks": remarks}).Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "状态更新失败")
		return
	}
	tx.Commit()
	response.Success(ctx, nil, "ok")
}
