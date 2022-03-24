package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/dto"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/service"
	"kuukaa.fun/danmu-v4/util"
)

/*********************************************************
** 函数功能: 发送注册验证码
** 日    期: 2021/7/23
**********************************************************/
func SendCode(ctx *gin.Context) {
	var requestUser dto.SendCodeDto
	err := ctx.Bind(&requestUser)
	if err != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}
	email := requestUser.Email

	//数据验证
	if !util.VerifyEmailFormat(email) {
		response.CheckFail(ctx, nil, response.EmailFormatCheck)
		return
	}

	res := service.SendRegisterCodeService(email)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 给自己发送验证码
** 说    明: 发送前需要验证邮箱与用户的邮箱是否一致
** 日    期:2021/10/25
**********************************************************/
func SendCodeToMyself(ctx *gin.Context) {
	var requestUser dto.SendCodeDto
	err := ctx.Bind(&requestUser)
	if err != nil {
		response.Response(ctx, http.StatusBadRequest, 4000, nil, response.RequestError)
		return
	}
	email := requestUser.Email
	if !util.VerifyEmailFormat(email) {
		response.CheckFail(ctx, nil, response.EmailFormatCheck)
		return
	}

	uid, _ := ctx.Get("id")
	res := service.SendCodeToMyselfService(email, uid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 发送登录验证码
** 日    期: 2022年2月10日12:55:40
**********************************************************/
func SendLoginCode(ctx *gin.Context) {
	var requestUser dto.SendCodeDto
	err := ctx.Bind(&requestUser)
	if err != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}
	email := requestUser.Email
	if !util.VerifyEmailFormat(email) {
		response.CheckFail(ctx, nil, response.EmailFormatCheck)
		return
	}

	res := service.SendLoginCodeService(email)
	response.HandleResponse(ctx, res)
}
