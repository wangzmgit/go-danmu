package controller

import (
	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/dto"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/service"
	"kuukaa.fun/danmu-v4/util"
)

/*********************************************************
** 函数功能: 创建反馈(站内)
** 日    期: 2021年12月3日15:02:02
** 版    本: 3.6.6
**********************************************************/
func CreateOpinionOnSite(ctx *gin.Context) {
	var request dto.OpinionOnSiteDto
	err := ctx.Bind(&request)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	desc := request.Desc

	uid, _ := ctx.Get("id")

	//验证数据
	if len(desc) == 0 {
		response.CheckFail(ctx, nil, "内容不能为空")
		return
	}

	res := service.CreateOpinionOnSiteService(desc, uid.(uint))
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 创建反馈
** 日    期: 2021年12月3日17:07:46
** 版    本: 3.6.6
**********************************************************/
func CreateOpinion(ctx *gin.Context) {
	var request dto.OpinionDto
	err := ctx.Bind(&request)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	name := request.Name
	email := request.Email
	telephone := request.Telephone
	gender := request.Gender
	desc := request.Desc

	//验证数据
	if len(name) == 0 {
		response.CheckFail(ctx, nil, "姓名不能为空")
		return
	}

	if len(email) != 0 && !util.VerifyEmailFormat(email) {
		response.CheckFail(ctx, nil, "邮箱格式有误")
		return
	}

	if len(telephone) != 0 && !util.VerifyTelephoneFormat(telephone) {
		response.CheckFail(ctx, nil, "联系方式格式有误")
		return
	}

	if gender < 0 || gender > 2 {
		response.CheckFail(ctx, nil, "性别选择有误")
		return
	}

	if len(desc) == 0 {
		response.CheckFail(ctx, nil, "内容不能为空")
		return
	}

	res := service.CreateOpinionService(request)
	response.HandleResponse(ctx, res)
}
