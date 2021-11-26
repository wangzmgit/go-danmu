package service

import (
	"net/http"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/vo"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

/*********************************************************
** 函数功能: 管理员登录
** 日    期: 2021年11月12日11:39:10
**********************************************************/
func AdminLoginService(email string, password string) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}
	var admin model.Admin

	DB := common.GetDB()
	DB.Where("email = ?", email).First(&admin)
	if admin.ID != 0 {
		var adminInfo vo.AdminVo
		//判断密码
		if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err == nil {
			token, err := common.ReleaseAdminToken(admin.ID)
			if err != nil {
				res.HttpStatus = http.StatusInternalServerError
				res.Code = response.ServerErrorCode
				res.Msg = "服务器出错了"
				return res
			}
			adminInfo.Name = admin.Name
			adminInfo.Authority = admin.Authority
			res.Data = gin.H{"token": token, "info": adminInfo}
			return res
		}
	}
	res.HttpStatus = http.StatusUnprocessableEntity
	res.Code = response.CheckFailCode
	res.Msg = "用户名或密码错误"
	return res
}

/*********************************************************
** 函数功能: 添加管理员
** 日    期: 2021年11月12日11:47:20
**********************************************************/
func AddAdminService(adminDto dto.AddAdminRequest) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}
	var admin model.Admin
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminDto.Password), bcrypt.DefaultCost)
	if err != nil {
		res.HttpStatus = http.StatusInternalServerError
		res.Code = response.ServerErrorCode
		res.Msg = "服务器出错了"
		return res
	}
	DB := common.GetDB()

	DB.Where("email = ?", adminDto.Email).First(&admin)
	if admin.ID == 0 {
		newAdmin := model.Admin{
			Name:      adminDto.Name,
			Email:     adminDto.Email,
			Password:  string(hashedPassword),
			Authority: adminDto.Authority,
		}
		DB.Create(&newAdmin)
		return res
	}
	res.HttpStatus = http.StatusUnprocessableEntity
	res.Code = response.CheckFailCode
	res.Msg = "该管理员账号已存在"
	return res
}

/*********************************************************
** 函数功能: 删除管理员
** 日    期: 2021年11月12日11:52:06
**********************************************************/
func DeleteAdminService(id uint) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}
	var admin model.Admin

	DB := common.GetDB()
	DB.First(&admin, id)
	if admin.ID != 0 {
		DB.Delete(&admin)
		return res
	}
	res.HttpStatus = http.StatusBadRequest
	res.Code = response.FailCode
	res.Msg = "该管理员账号不存在"
	return res
}

/*********************************************************
** 函数功能: 获取管理员列表
** 日    期: 2021年11月12日11:55:39
**********************************************************/
func GetAdminListService(page int, pageSize int) response.ResponseStruct {
	var total int //记录总数
	var admins []vo.AdminListVo

	DB := common.GetDB()
	DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
	DB.Model(&model.Admin{}).Select("id,name,email,authority").Scan(&admins).Count(&total)

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"count": total, "admins": admins},
		Msg:        "ok",
	}
}
