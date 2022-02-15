package manage

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"kuukaa.fun/danmu-v4/common"
	"kuukaa.fun/danmu-v4/dto"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/service"
	"kuukaa.fun/danmu-v4/util"
	"kuukaa.fun/danmu-v4/vo"
)

/*********************************************************
** 函数功能: 管理员登录
** 日    期:2021/7/10
**********************************************************/
func AdminLogin(ctx *gin.Context) {
	//获取参数
	var requestAdmin dto.AdminLoginDto
	requestErr := ctx.Bind(&requestAdmin)
	if requestErr != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}
	email := requestAdmin.Email
	password := requestAdmin.Password

	//数据验证
	if len(email) == 0 || len(password) == 0 {
		response.CheckFail(ctx, nil, response.LoginCheck)
		return
	}
	//默认管理员
	if email == viper.GetString("admin.email") && password == viper.GetString("admin.password") {
		//发放token
		var adminInfo vo.AdminVo
		token, err := common.ReleaseAdminToken(0)
		if err != nil {
			response.ServerError(ctx, nil, response.SystemError)
			return
		}
		adminInfo.Name = "管理员"
		adminInfo.Authority = util.SuperAdmin
		response.Success(ctx, gin.H{"token": token, "info": adminInfo}, response.OK)
	} else {
		//查询管理员表
		res := service.AdminLoginService(email, password)
		response.HandleResponse(ctx, res)
	}
}

/*********************************************************
** 函数功能: 添加管理员
** 日    期:2021/8/1
**********************************************************/
func AddAdmin(ctx *gin.Context) {
	//获取参数
	var request dto.AddAdminDto
	err := ctx.Bind(&request)
	if err != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}
	name := request.Name
	email := request.Email
	password := request.Password
	authority := request.Authority

	//数据验证
	if len(name) == 0 {
		response.CheckFail(ctx, nil, response.NameCheck)
		return
	}
	if !util.VerifyEmailFormat(email) {
		response.CheckFail(ctx, nil, response.EmailFormatCheck)
		return
	}
	if len(password) < 6 {
		response.CheckFail(ctx, nil, response.PasswordCheck)
		return
	}
	if authority != util.Admin && authority != util.Auditor {
		response.CheckFail(ctx, nil, response.AuthorityCheck)
		return
	}

	res := service.AddAdminService(request)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 删除管理员
** 日    期: 2021/8/3
**********************************************************/
func DeleteAdmin(ctx *gin.Context) {
	//获取参数
	var request dto.AdminIdDto
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}
	id := request.ID
	if id == 0 {
		response.CheckFail(ctx, nil, response.UserNotExist)
		return
	}
	res := service.DeleteAdminService(id)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 获取管理员列表
** 日    期:2021/8/2
**********************************************************/
func GetAdminList(ctx *gin.Context) {
	//获取参数
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	if page <= 0 || pageSize <= 0 {
		response.CheckFail(ctx, nil, response.PageOrSizeError)
		return
	}

	res := service.GetAdminListService(page, pageSize)
	response.HandleResponse(ctx, res)
}
