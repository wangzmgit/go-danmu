package admin_controller

import (
	"strconv"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/service"
	"wzm/danmu3.0/util"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

/*********************************************************
** 函数功能: 管理员登录
** 日    期:2021/7/10
**********************************************************/
func AdminLogin(ctx *gin.Context) {
	//获取参数
	var requestAdmin dto.AdminLoginRequest
	requestErr := ctx.Bind(&requestAdmin)
	if requestErr != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	email := requestAdmin.Email
	password := requestAdmin.Password

	//数据验证
	if len(email) == 0 || len(password) == 0 {
		response.CheckFail(ctx, nil, "请输入用户名或密码")
		return
	}
	//默认管理员
	if email == viper.GetString("admin.email") && password == viper.GetString("admin.password") {
		//发放token
		token, err := common.ReleaseAdminToken(0)
		if err != nil {
			response.ServerError(ctx, nil, "系统异常")
			return
		}
		response.Success(ctx, gin.H{"token": token}, "ok")
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
	var request dto.AddAdminRequest
	err := ctx.Bind(&request)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	email := request.Email
	password := request.Password
	authority := request.Authority

	//数据验证
	if !util.VerifyEmailFormat(email) {
		response.CheckFail(ctx, nil, "邮箱格式有误")
		return
	}
	if len(password) < 6 {
		response.CheckFail(ctx, nil, "密码不要少于六位")
		return
	}
	if authority != util.Admin && authority != util.Auditor {
		response.CheckFail(ctx, nil, "权限选择有误")
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
	var request dto.AdminIDRequest
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	id := request.ID
	if id == 0 {
		response.CheckFail(ctx, nil, "该管理员账号不存在")
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
		response.CheckFail(ctx, nil, "页码或数量有误")
		return
	}

	res := service.GetAdminListService(page, pageSize)
	response.HandleResponse(ctx, res)
}
