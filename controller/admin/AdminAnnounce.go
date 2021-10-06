package admin_controller

import (
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"

	"github.com/gin-gonic/gin"
)

/*********************************************************
** 函数功能: 获取公告
** 日    期:2021/8/4
**********************************************************/
func AdminGetAnnounce(ctx *gin.Context) {
	DB := common.GetDB()
	var announceList []dto.AdminAnnounceDto
	DB.Raw("select id,created_at,title,content,url from announces where deleted_at is null").Scan(&announceList)
	response.Success(ctx, gin.H{"announces": announceList}, "ok")
}

/*********************************************************
** 函数功能: 添加公告
** 日    期:2021/8/4
**********************************************************/
func AddAnnounce(ctx *gin.Context) {
	DB := common.GetDB()
	var announce = model.Announce{}
	err := ctx.Bind(&announce)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	title := announce.Title
	content := announce.Content
	url := announce.Url
	if len(title) == 0 {
		response.CheckFail(ctx, nil, "标题不能为空")
		return
	}
	if len(content) == 0 {
		response.CheckFail(ctx, nil, "内容不能为空")
		return
	}
	newAnnounce := model.Announce{
		Title:   title,
		Content: content,
		Url:     url,
	}
	DB.Create(&newAnnounce)
	//返回结果
	response.Success(ctx, nil, "ok")
}

/*********************************************************
** 函数功能: 删除公告
** 日    期:2021/8/4
**********************************************************/
func DeleteAnnounce(ctx *gin.Context) {
	DB := common.GetDB()
	var request = AdminIDRequest{}
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	id := request.ID
	DB.Where("id = ?", id).Delete(model.Announce{})
	response.Success(ctx, nil, "ok")
}
