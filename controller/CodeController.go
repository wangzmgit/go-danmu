package controller

import (
	"net/http"
	"time"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/util"

	"github.com/gin-gonic/gin"
)

/*********************************************************
** 函数功能: 发送验证码
** 日    期:2021/7/23
**********************************************************/
func SendCode(ctx *gin.Context) {
	var requestUser = model.User{}
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
	//邮箱是否存在
	DB := common.GetDB()
	if IsEmailExist(DB, email) {
		response.CheckFail(ctx, nil, "该邮箱已经被注册了")
		return
	}
	//存储code到redis
	code, _ := common.RedisClient.Get(util.CodeKey(email)).Result()
	if code != "" {
		//如果时间小于一分钟则不能重新发送
		duration, _ := common.RedisClient.TTL(util.CodeKey(email)).Result()
		if duration >= 240000000000 {
			response.Fail(ctx, nil, "操作过于频繁")
			return
		}
	}
	randomCode := util.RandomCode(6)
	err = common.RedisClient.Set(util.CodeKey(email), randomCode, time.Second*300).Err()
	if err != nil {
		response.ServerError(ctx, nil, "发送失败")
		return
	}
	send := util.SendEmail(email, randomCode)
	if send {
		response.Success(ctx, nil, "发送成功")
	} else {
		common.RedisClient.Del(util.CodeKey(email))
		response.Fail(ctx, nil, "发送失败")
	}
}

/*********************************************************
** 函数功能: 验证验证码
** 日    期:2021/7/24
**********************************************************/
func VerificationCode(email string, code string) bool {
	dbCode, _ := common.RedisClient.Get(util.CodeKey(email)).Result()
	if dbCode == "" || dbCode != code {
		return false
	}
	return true
}

/*********************************************************
** 函数功能: 给自己发送验证码
** 说    明: 发送前需要验证邮箱与用户的邮箱是否一致
** 日    期:2021/10/25
**********************************************************/
func SendCodeToMyself(ctx *gin.Context) {
	var requestUser = model.User{}
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
	var user model.User
	DB := common.GetDB()
	uid, _ := ctx.Get("id")
	DB.First(&user, uid)
	if user.Email != email {
		response.CheckFail(ctx, nil, "邮箱验证失败")
		return
	}
	//存储code到redis
	code, _ := common.RedisClient.Get(util.CodeKey(email)).Result()
	if code != "" {
		//如果时间小于一分钟则不能重新发送
		duration, _ := common.RedisClient.TTL(util.CodeKey(email)).Result()
		if duration >= 240000000000 {
			response.Fail(ctx, nil, "操作过于频繁")
			return
		}
	}
	randomCode := util.RandomCode(6)
	err = common.RedisClient.Set(util.CodeKey(email), randomCode, time.Second*300).Err()
	if err != nil {
		response.ServerError(ctx, nil, "发送失败")
		return
	}
	send := util.SendEmail(email, randomCode)
	if send {
		response.Success(ctx, nil, "ok")
	} else {
		common.RedisClient.Del(util.CodeKey(email))
		response.Fail(ctx, nil, "发送失败")
	}
}
