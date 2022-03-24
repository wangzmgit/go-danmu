package service

import (
	"net/http"
	"time"

	"kuukaa.fun/danmu-v4/common"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/util"
)

/*********************************************************
** 函数功能: 发送注册验证码
** 日    期: 2022年3月24日12:18:40
**********************************************************/
func SendRegisterCodeService(email string) response.ResponseStruct {
	//邮箱是否存在
	if isEmailRegistered(email) {
		return response.ResponseStruct{
			HttpStatus: http.StatusUnprocessableEntity,
			Code:       response.CheckFailCode,
			Data:       nil,
			Msg:        response.EmailRegistered,
		}
	}

	return sendCode(email, util.CodeKey(email), util.RegisterCode)
}

/*********************************************************
** 函数功能: 给自己发送验证码
** 日    期: 2022年3月24日12:29:34
**********************************************************/
func SendCodeToMyselfService(email string, uid interface{}) response.ResponseStruct {
	//邮箱是否属于当前用户
	if !isEmailBelongsToCurrentUser(email, uid) {
		return response.ResponseStruct{
			HttpStatus: http.StatusUnprocessableEntity,
			Code:       response.CheckFailCode,
			Data:       nil,
			Msg:        response.VerificationFail,
		}
	}

	return sendCode(email, util.CodeKey(email), util.ModifyPasswordCode)
}

/*********************************************************
** 函数功能: 给自己发送验证码
** 日    期: 2022年3月24日12:37:24
**********************************************************/
func SendLoginCodeService(email string) response.ResponseStruct {
	return sendCode(email, util.LoginCodeKey(email), util.LoginCode)
}

/*********************************************************
** 函数功能: 验证验证码
** 日    期:2021/7/24
**********************************************************/
func VerificationCode(emailKey string, code string) bool {
	if len(code) == 0 {
		return false
	}
	Redis := common.RedisClient
	if Redis == nil {
		util.Logfile(util.ErrorLog, "Verification code redis error")
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
** 函数功能: 发送验证码
** 日    期: 2022年3月24日12:19:28
**********************************************************/
func sendCode(email, key, subject string) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}
	//存储code到redis
	Redis := common.RedisClient
	if Redis == nil {
		res.HttpStatus = http.StatusInternalServerError
		res.Code = response.ServerErrorCode
		res.Msg = response.SystemError
		return res
	}
	code, _ := Redis.Get(key).Result()
	if code != "" {
		//如果时间小于一分钟则不能重新发送
		duration, _ := Redis.TTL(key).Result()
		if duration >= 240000000000 {
			res.HttpStatus = http.StatusBadRequest
			res.Code = response.FailCode
			res.Msg = response.OperationTooFrequently
			return res
		}
	}
	randomCode := util.RandomCode(6)
	err := Redis.Set(key, randomCode, time.Second*300).Err()
	if err != nil {
		res.HttpStatus = http.StatusInternalServerError
		res.Code = response.ServerErrorCode
		res.Msg = response.SendFail
		return res
	}
	send := util.SendEmail(email, randomCode, subject)
	if send {
		return res
	} else {
		Redis.Del(key)
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = response.SendFail
		return res
	}
}
