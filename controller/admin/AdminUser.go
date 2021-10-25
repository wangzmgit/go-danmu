package admin_controller

import (
	"strconv"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/util"

	"github.com/gin-gonic/gin"
)

/*********************************************************
** 函数功能: 获取用户列表
** 日    期:2021/8/3
**********************************************************/
func GetUserList(ctx *gin.Context) {
	DB := common.GetDB()
	var users []dto.AdminUserDto
	//获取分页信息
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	if page > 0 && pageSize > 0 {
		//记录总数
		var total int
		DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
		//评论
		DB.Model(&model.User{}).Select("id,name,created_at,email,avatar,sign,gender").Scan(&users).Count(&total)
		response.Success(ctx, gin.H{"count": total, "users": users}, "ok")
	} else {
		response.Fail(ctx, nil, "获取数量有误")
	}
}

/*********************************************************
** 函数功能: 修改用户信息
** 日    期:2021/8/3
**********************************************************/
func AdminModifyUser(ctx *gin.Context) {
	DB := common.GetDB()
	var requestUser model.User
	var user model.User
	err := ctx.Bind(&requestUser)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	id := requestUser.ID
	email := requestUser.Email
	name := requestUser.Name
	sign := requestUser.Sign
	if len(name) == 0 {
		response.CheckFail(ctx, nil, "昵称不能为空")
		return
	}
	if !util.VerifyEmailFormat(email) {
		response.CheckFail(ctx, nil, "邮箱格式有误")
		return
	}
	//新邮箱的uid不为当前uid
	DB.Where("email = ? and id <> ?", email, id).First(&user)
	if user.ID != 0 {
		response.CheckFail(ctx, nil, "邮箱已存在")
		return
	}
	DB.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{"email": email, "name": name, "sign": sign})
	response.Success(ctx, nil, "ok")
}

/*********************************************************
** 函数功能: 删除用户
** 日    期:2021/8/3
**********************************************************/
func AdminDeleteUser(ctx *gin.Context) {
	DB := common.GetDB()
	var request = AdminIDRequest{}
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	id := request.ID
	DB.Where("id = ?", id).Delete(model.User{})
	response.Success(ctx, nil, "ok")
}
