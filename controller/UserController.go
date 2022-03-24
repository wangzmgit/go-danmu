package controller

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"kuukaa.fun/danmu-v4/dto"
	"kuukaa.fun/danmu-v4/model"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/service"
	"kuukaa.fun/danmu-v4/util"
	"kuukaa.fun/danmu-v4/vo"
)

/*********************************************************
** 函数功能: 用户注册
** 日    期: 2021/7/10
**********************************************************/
func Register(ctx *gin.Context) {
	//获取参数
	var request dto.RegisterDto
	err := ctx.Bind(&request)
	if err != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}
	name := request.Name
	email := request.Email
	password := request.Password
	code := request.Code

	//数据验证
	if !util.VerifyEmailFormat(email) {
		response.CheckFail(ctx, nil, response.EmailFormatCheck)
		return
	}
	if len(password) < 6 {
		response.CheckFail(ctx, nil, response.PasswordCheck)
		return
	}
	if !service.VerificationCode(util.CodeKey(email), code) {
		response.CheckFail(ctx, nil, response.VerificationCodeError)
		return
	}
	//如果名称为空，则为随机字符串
	if len(name) == 0 {
		name = util.RandomString(10)
	}
	res := service.RegisterService(request)
	//返回结果
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 用户登录
** 日    期:2021/7/10
**********************************************************/
func Login(ctx *gin.Context) {
	//获取参数
	var request dto.LoginDto
	err := ctx.Bind(&request)
	if err != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}
	email := request.Email
	password := request.Password

	//数据验证
	if !util.VerifyEmailFormat(email) {
		response.CheckFail(ctx, nil, response.EmailFormatCheck)
		return
	}
	if len(password) < 6 {
		response.CheckFail(ctx, nil, response.PasswordCheck)
		return
	}
	res := service.LoginService(request, ctx.ClientIP())

	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 邮箱验证登录
** 日    期: 2022年2月10日13:05:59
**********************************************************/
func EmailLogin(ctx *gin.Context) {
	//获取参数
	var request dto.EmailLoginDto
	err := ctx.Bind(&request)
	if err != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}
	email := request.Email
	code := request.Code

	//数据验证
	if !util.VerifyEmailFormat(email) {
		response.CheckFail(ctx, nil, response.EmailFormatCheck)
		return
	}

	//验证验证码是否正确
	if !service.VerificationCode(util.LoginCodeKey(email), code) {
		response.CheckFail(ctx, nil, response.VerificationCodeError)
		return
	}

	res := service.EmailLoginService(request, ctx.ClientIP())
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 用户获取个人信息
** 日    期:2021/7/10
**********************************************************/
func UserInfo(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	response.Success(ctx, gin.H{"data": vo.ToUserVo(user.(model.User))}, response.OK)
}

/*********************************************************
** 函数功能: 用户修改个人信息
** 日    期:2021/7/10
**********************************************************/
func ModifyInfo(ctx *gin.Context) {
	var request dto.ModifyUserDto
	err := ctx.Bind(&request)
	if err != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}
	name := request.Name
	birthday := request.Birthday

	if len(name) == 0 {
		response.CheckFail(ctx, nil, response.NickCheck)
		return
	}
	//判断日期
	tBirthday, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		response.CheckFail(ctx, nil, response.BirthdayFormatCheck)
		return
	}

	//从上下文中获取用户id
	id, _ := ctx.Get("id")
	res := service.UserModifyService(request, id, tBirthday)

	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 用户修改密码
** 日    期:2021/7/10
**********************************************************/
func ModifyPassword(ctx *gin.Context) {
	//获取参数
	var passModify dto.ModifyPasswordDto
	err := ctx.Bind(&passModify)
	if err != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}
	password := passModify.Password
	code := passModify.Code

	//数据验证
	if len(password) < 6 {
		response.CheckFail(ctx, nil, response.PasswordCheck)
		return
	}
	//从上下文中获取用户信息
	user, _ := ctx.Get("user")
	modelUser := user.(model.User)
	//验证验证码是否正确
	if !service.VerificationCode(util.CodeKey(modelUser.Email), code) {
		response.CheckFail(ctx, nil, response.VerificationCodeError)
		return
	}

	//加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.CheckFail(ctx, nil, response.SystemError)
		//记录日志
		util.Logfile(util.ErrorLog, " hashed password "+err.Error())
		return
	}

	//更新数据
	res := service.ModifyPasswordService(string(hashedPassword), modelUser)
	response.HandleResponse(ctx, res)

}

/*********************************************************
** 函数功能: 通过用户ID获取用户信息
** 日    期: 2021/7/10
** 说    明: 用于获取其他用户信息
**********************************************************/
func GetUserInfoByID(ctx *gin.Context) {
	uid, _ := strconv.Atoi(ctx.Query("uid"))
	res := service.GetUserInfoByIDService(uid)
	response.HandleResponse(ctx, res)
}
