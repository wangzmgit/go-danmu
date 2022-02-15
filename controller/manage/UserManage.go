package manage

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/dto"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/service"
	"kuukaa.fun/danmu-v4/util"
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
		response.CheckFail(ctx, nil, response.PageOrSizeError)
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
	var requestUser dto.AdminModifyUserDto
	err := ctx.Bind(&requestUser)
	if err != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}
	email := requestUser.Email
	name := requestUser.Name

	if len(name) == 0 {
		response.CheckFail(ctx, nil, response.NickCheck)
		return
	}
	if !util.VerifyEmailFormat(email) {
		response.CheckFail(ctx, nil, response.EmailFormatCheck)
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
	var request dto.AdminIdDto
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}
	id := request.ID

	res := service.AdminDeleteUserService(id)
	response.HandleResponse(ctx, res)
}
