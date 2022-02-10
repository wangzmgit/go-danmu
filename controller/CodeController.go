package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/common"
	"kuukaa.fun/danmu-v4/dto"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/service"
	"kuukaa.fun/danmu-v4/util"
)

/*********************************************************
** 函数功能: 发送验证码
** 日    期:2021/7/23
**********************************************************/
func SendCode(ctx *gin.Context) {
	var requestUser dto.SendCodeDto
	err := ctx.Bind(&requestUser)
	if err != nil {
		response.Response(ctx, http.StatusBadRequest, 4000, nil, "请求错误")
		return
	}
	email := requestUser.Email

	//数据验证
	if !util.VerifyEmailFormat(email) {
		response.CheckFail(ctx, nil, "邮箱格式有误")
		return
	}
	//邮箱是否存在
	if service.IsEmailRegistered(email) {
		response.CheckFail(ctx, nil, "该邮箱已经被注册了")
		return
	}
	//存储code到redis
	Redis := common.RedisClient
	if Redis == nil {
		response.ServerError(ctx, nil, "系统故障")
		return
	}
	code, _ := Redis.Get(util.CodeKey(email)).Result()
	if code != "" {
		//如果时间小于一分钟则不能重新发送
		duration, _ := Redis.TTL(util.CodeKey(email)).Result()
		if duration >= 240000000000 {
			response.Fail(ctx, nil, "操作过于频繁")
			return
		}
	}
	randomCode := util.RandomCode(6)
	err = Redis.Set(util.CodeKey(email), randomCode, time.Second*300).Err()
	if err != nil {
		response.ServerError(ctx, nil, "发送失败")
		return
	}
	send := util.SendEmail(email, randomCode, "注册验证码")
	if send {
		response.Success(ctx, nil, "发送成功")
	} else {
		Redis.Del(util.CodeKey(email))
		response.Fail(ctx, nil, "发送失败")
	}
}

/*********************************************************
** 函数功能: 验证验证码
** 日    期:2021/7/24
**********************************************************/
func VerificationCode(emailKey string, code string) bool {
	Redis := common.RedisClient
	if Redis == nil {
		util.Logfile("[Error]", "Verification code redis error")
		return false
	}
	dbCode, _ := Redis.Get(emailKey).Result()
	if dbCode == "" || dbCode != code {
		return false
	}
	Redis.Del(emailKey)
	return true
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
		response.Response(ctx, http.StatusBadRequest, 4000, nil, "请求错误")
		return
	}
	email := requestUser.Email
	if !util.VerifyEmailFormat(email) {
		response.CheckFail(ctx, nil, "邮箱格式有误")
		return
	}
	//邮箱是否属于当前用户
	uid, _ := ctx.Get("id")
	if !service.IsEmailBelongsToCurrentUser(email, uid) {
		response.CheckFail(ctx, nil, "邮箱验证失败")
		return
	}
	//存储code到redis
	Redis := common.RedisClient
	if Redis == nil {
		response.ServerError(ctx, nil, "系统故障")
		return
	}
	code, _ := Redis.Get(util.CodeKey(email)).Result()
	if code != "" {
		//如果时间小于一分钟则不能重新发送
		duration, _ := Redis.TTL(util.CodeKey(email)).Result()
		if duration >= 240000000000 {
			response.Fail(ctx, nil, "操作过于频繁")
			return
		}
	}
	randomCode := util.RandomCode(6)
	err = Redis.Set(util.CodeKey(email), randomCode, time.Second*300).Err()
	if err != nil {
		response.ServerError(ctx, nil, "发送失败")
		return
	}
	send := util.SendEmail(email, randomCode, "修改密码验证")
	if send {
		response.Success(ctx, nil, "ok")
	} else {
		Redis.Del(util.CodeKey(email))
		response.Fail(ctx, nil, "发送失败")
	}
}

/*********************************************************
** 函数功能: 发送登录验证码
** 日    期: 2022年2月10日12:55:40
**********************************************************/
func SendLoginCode(ctx *gin.Context) {
	var requestUser dto.SendCodeDto
	err := ctx.Bind(&requestUser)
	if err != nil {
		response.Response(ctx, http.StatusBadRequest, 4000, nil, "请求错误")
		return
	}
	email := requestUser.Email
	if !util.VerifyEmailFormat(email) {
		response.CheckFail(ctx, nil, "邮箱格式有误")
		return
	}
	//存储code到redis
	Redis := common.RedisClient
	if Redis == nil {
		response.ServerError(ctx, nil, "系统故障")
		return
	}
	code, _ := Redis.Get(util.LoginCodeKey(email)).Result()
	if code != "" {
		//如果时间小于一分钟则不能重新发送
		duration, _ := Redis.TTL(util.LoginCodeKey(email)).Result()
		if duration >= 240000000000 {
			response.Fail(ctx, nil, "操作过于频繁")
			return
		}
	}
	randomCode := util.RandomCode(6)
	err = Redis.Set(util.LoginCodeKey(email), randomCode, time.Second*300).Err()
	if err != nil {
		response.ServerError(ctx, nil, "发送失败")
		return
	}
	send := util.SendEmail(email, randomCode, "登录验证")
	if send {
		response.Success(ctx, nil, "ok")
	} else {
		Redis.Del(util.LoginCodeKey(email))
		response.Fail(ctx, nil, "发送失败")
	}
}
