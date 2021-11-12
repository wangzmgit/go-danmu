package admin_controller

import (
	"strconv"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/service"
	"wzm/danmu3.0/util"

	"github.com/gin-gonic/gin"
)

/*********************************************************
** 函数功能: 获取用户列表
** 日    期: 2021/8/3
**********************************************************/
func GetUserList(ctx *gin.Context) {
	//获取分页信息
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	if page <= 0 || pageSize <= 0 {
		response.CheckFail(ctx, nil, "页码或数量有误")
		return
	}

	res := service.GetUserListService(page, pageSize)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 修改用户信息
** 日    期: 2021/8/3
**********************************************************/
func AdminModifyUser(ctx *gin.Context) {
	//获取参数
	var requestUser dto.AdminModifyUserRequest
	err := ctx.Bind(&requestUser)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	email := requestUser.Email
	name := requestUser.Name

	if len(name) == 0 {
		response.CheckFail(ctx, nil, "昵称不能为空")
		return
	}
	if !util.VerifyEmailFormat(email) {
		response.CheckFail(ctx, nil, "邮箱格式有误")
		return
	}

	res := service.AdminModifyUserService(requestUser)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 删除用户
** 日    期: 2021/8/3
**********************************************************/
func AdminDeleteUser(ctx *gin.Context) {
	var request dto.AdminIDRequest
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	id := request.ID

	res := service.AdminDeleteUserService(id)
	response.HandleResponse(ctx, res)
}
