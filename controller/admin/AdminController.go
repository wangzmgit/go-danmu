package admin_controller

import (
	"strconv"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/util"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

const (
	SuperAdmin = 3000
	Admin      = 2000
	Auditor    = 1000
)

type AdminIDRequest struct {
	ID uint
}

/*********************************************************
** 函数功能: 管理员登录
** 日    期:2021/7/10
**********************************************************/
func AdminLogin(ctx *gin.Context) {
	//获取参数
	type requestLogin struct {
		Email    string
		Password string
	}
	var requestAdmin requestLogin
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
		DB := common.GetDB()
		var admin model.Admin
		DB.Where("email = ?", email).First(&admin)
		if admin.ID != 0 {
			//判断密码
			if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err == nil {
				token, err := common.ReleaseAdminToken(admin.ID)
				if err != nil {
					response.ServerError(ctx, nil, "系统异常")
					return
				}
				response.Success(ctx, gin.H{"token": token}, "ok")
				return
			}
		}
		response.CheckFail(ctx, nil, "用户名或密码错误！")
	}
}

/*********************************************************
** 函数功能: 添加管理员
** 日    期:2021/8/1
**********************************************************/
func AddAdmin(ctx *gin.Context) {
	DB := common.GetDB()
	var request model.Admin
	var admin model.Admin
	err := ctx.Bind(&request)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	name := request.Name
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

	if authority != Admin && authority != Auditor {
		response.CheckFail(ctx, nil, "权限选择有误")
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	DB.Where("email = ?", email).First(&admin)
	if admin.ID == 0 {
		newAdmin := model.Admin{
			Name:      name,
			Email:     email,
			Password:  string(hashedPassword),
			Authority: authority,
		}
		DB.Create(&newAdmin)
		response.Success(ctx, nil, "ok")
	} else {
		response.CheckFail(ctx, nil, "用户已存在")
	}
}

/*********************************************************
** 函数功能: 删除管理员
** 日    期:2021/8/3
**********************************************************/
func DeleteAdmin(ctx *gin.Context) {
	DB := common.GetDB()
	var request = AdminIDRequest{}
	var admin = model.Admin{}
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	id := request.ID
	DB.First(&admin, id)
	if admin.ID != 0 {
		DB.Delete(&admin)
		response.Success(ctx, nil, "ok")
	} else {
		response.Fail(ctx, nil, "管理员不存在")
	}
}

/*********************************************************
** 函数功能: 获取管理员列表
** 日    期:2021/8/2
**********************************************************/
func GetAdminList(ctx *gin.Context) {
	DB := common.GetDB()
	var admins []dto.AdminListDto
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	if page > 0 && pageSize > 0 {
		//记录总数
		var total int
		DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
		//评论
		DB.Model(&model.Admin{}).Select("id,name,email,authority").Scan(&admins).Count(&total)
		response.Success(ctx, gin.H{"count": total, "admins": admins}, "ok")
	} else {
		response.Fail(ctx, nil, "获取数量有误")
	}
}
