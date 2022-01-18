package controller

import (
	"net/http"
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
		response.Fail(ctx, nil, "请求错误")
		return
	}
	name := request.Name
	email := request.Email
	password := request.Password
	code := request.Code

	//数据验证
	if !util.VerifyEmailFormat(email) {
		response.CheckFail(ctx, nil, "邮箱格式有误哦")
		return
	}
	if len(password) < 6 {
		response.CheckFail(ctx, nil, "密码不要少于六位")
		return
	}
	if !VerificationCode(email, code) {
		response.CheckFail(ctx, nil, "验证码有误")
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
		response.Response(ctx, http.StatusBadRequest, 4000, nil, "请求错误")
		return
	}
	email := request.Email
	password := request.Password

	//数据验证
	if !util.VerifyEmailFormat(email) {
		response.CheckFail(ctx, nil, "邮箱格式有误哦")
		return
	}
	if len(password) < 6 {
		response.CheckFail(ctx, nil, "密码不要少于六位")
		return
	}
	res := service.LoginService(request, ctx.ClientIP())

	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 用户获取个人信息
** 日    期:2021/7/10
**********************************************************/
func UserInfo(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	response.Success(ctx, gin.H{"data": vo.ToUserVo(user.(model.User))}, "ok")
}

/*********************************************************
** 函数功能: 用户修改个人信息
** 日    期:2021/7/10
**********************************************************/
func ModifyInfo(ctx *gin.Context) {
	//获取参数
	var request dto.ModifyUserDto
	err := ctx.Bind(&request)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	name := request.Name
	birthday := request.Birthday

	if len(name) == 0 {
		response.CheckFail(ctx, nil, "昵称不能为空哦")
		return
	}
	//判断日期
	tBirthday, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		response.CheckFail(ctx, nil, "请输入正确的出生日期哦")
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
		response.Fail(ctx, nil, "请求错误")
		return
	}
	password := passModify.Password
	code := passModify.Code

	//数据验证
	if len(password) < 6 {
		response.CheckFail(ctx, nil, "密码不要少于六位")
		return
	}
	//从上下文中获取用户信息
	user, _ := ctx.Get("user")
	modelUser := user.(model.User)
	//验证验证码是否正确
	if !VerificationCode(modelUser.Email, code) {
		response.CheckFail(ctx, nil, "验证码有误")
		return
	}

	//加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.CheckFail(ctx, nil, "服务器出错了")
		//记录日志
		util.Logfile("[Error]", " hashed password "+err.Error())
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
