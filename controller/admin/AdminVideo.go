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
** 函数功能: 获取视频列表
** 日    期:2021/8/4
**********************************************************/
func AdminGetVideoList(ctx *gin.Context) {
	DB := common.GetDB()
	var videos []model.Video
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	if page > 0 && pageSize > 0 {
		//记录总数
		var total int
		DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
		DB.Where("review = 1").Find(&videos).Count(&total)
		response.Success(ctx, gin.H{"count": total, "videos": dto.ToAdminVideoDto(videos)}, "ok")
	} else {
		response.Fail(ctx, nil, "获取数量有误")
	}
}

/*********************************************************
** 函数功能: 删除视频
** 日    期:2021/8/3
**********************************************************/
func AdminDeleteVideo(ctx *gin.Context) {
	DB := common.GetDB()
	var request = AdminIDRequest{}
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	id := request.ID
	DB.Where("id = ?", id).Delete(model.Video{})
	response.Success(ctx, nil, "ok")
}

/*********************************************************
** 函数功能: 管理员导入视频
** 日    期:2021/10/6
**********************************************************/
func ImportVideo(ctx *gin.Context) {
	var request model.Video
	err := ctx.Bind(&request)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	title := request.Title
	cover := request.Cover
	introduction := request.Introduction
	video := request.Video
	//验证数据
	if len(title) == 0 {
		response.CheckFail(ctx, nil, "标题不能为空")
		return
	}
	if len(cover) == 0 {
		response.CheckFail(ctx, nil, "封面图不能为空")
		return
	}
	if len(video) == 0 {
		response.CheckFail(ctx, nil, "视频链接不能为空")
		return
	}
	newVideo := model.Video{
		Title:        title,
		Cover:        cover,
		Introduction: introduction,
		Original:     true,
		Uid:          0,
		VideoType:    "mp4",
		Video:        video,
		Review:       true,
	}
	DB := common.GetDB()
	if err := DB.Create(&newVideo).Error; err != nil {
		response.Fail(ctx, nil, "上传失败")
		return
	}
	response.Success(ctx, gin.H{"vid": newVideo.ID}, "ok")
}
