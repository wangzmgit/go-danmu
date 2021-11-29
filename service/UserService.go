package service

import (
	"net/http"
	"strconv"
	"time"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/util"
	"wzm/danmu3.0/vo"

	"github.com/gin-gonic/gin"

	"github.com/jinzhu/gorm"

	"golang.org/x/crypto/bcrypt"
)

/*********************************************************
** 函数功能: 注册
** 日    期:2021/11/8
**********************************************************/
func RegisterService(user dto.RegisterRequest) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}

	//邮箱是否存在
	DB := common.GetDB()
	if IsEmailExist(DB, user.Email) {
		res.HttpStatus = http.StatusUnprocessableEntity
		res.Code = response.CheckFailCode
		res.Msg = "该邮箱已经被注册了"
		return res
	}

	//创建用户
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		res.HttpStatus = http.StatusInternalServerError
		res.Code = response.ServerErrorCode
		res.Msg = "服务器出错了"
		//记录日志
		util.Logfile("[Error]", " hashed password "+err.Error())
		return res
	}

	newUser := model.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: string(hashedPassword),
	}
	DB.Create(&newUser)
	return res
}

/*********************************************************
** 函数功能: 登录
** 日    期:2021/11/8
**********************************************************/
func LoginService(login dto.LoginRequest, userIP string) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}

	//判断邮箱是否存在
	var user model.User
	DB := common.GetDB()
	DB.Where("email = ?", login.Email).First(&user)
	if user.ID == 0 {
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = "用户不存在"
		return res
	}
	//判断密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		res.HttpStatus = http.StatusUnprocessableEntity
		res.Code = response.CheckFailCode
		res.Msg = "用户名或密码错误"
		return res
	}
	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		res.HttpStatus = http.StatusInternalServerError
		res.Code = response.ServerErrorCode
		res.Msg = "系统异常"
		util.Logfile("[Error]", " token generate error  "+err.Error())
		return res
	}
	util.Logfile("[Info]", " Token issued successfully uid "+strconv.Itoa(int(user.ID))+" | "+userIP)
	//返回数据
	res.Data = gin.H{"token": token, "user": vo.ToUserVo(user)}
	return res
}

/*********************************************************
** 函数功能: 修改用户信息
** 日    期:2021/11/8
**********************************************************/
func UserModifyService(modify dto.UserModifyRequest, uid interface{}, tBirthday time.Time) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}

	DB := common.GetDB()
	err := DB.Model(model.User{}).Where("id = ?", uid).Updates(
		map[string]interface{}{"name": modify.Name, "gender": modify.Gender, "birthday": tBirthday, "sign": modify.Sign},
	).Error
	if err != nil {
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = "修改失败"
		return res
	}
	return res
}

/*********************************************************
** 函数功能: 修改密码
** 日    期:2021/11/10
**********************************************************/
func ModifyPasswordService(password string, user model.User) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}
	DB := common.GetDB()

	//更新密码
	err := DB.Model(&user).Update("password", password).Error
	if err != nil {
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = "修改失败"
		return res
	}

	return res
}

/*********************************************************
** 函数功能: 通过用户ID获取用户信息
** 日    期:2021/11/10
**********************************************************/
func GetUserInfoByIDService(uid interface{}) response.ResponseStruct {
	var user model.User
	DB := common.GetDB()
	DB.Select("id,name,sign,avatar,gender").Where("id = ?", uid).First(&user)
	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"user": vo.ToUserVo(user)},
		Msg:        "ok",
	}
}

/*********************************************************
** 函数功能: 管理员获取用户列表
** 日    期: 2021年11月12日15:13:53
**********************************************************/
func GetUserListService(page int, pageSize int) response.ResponseStruct {
	var users []vo.AdminUserVo
	DB := common.GetDB()
	//记录总数
	var total int
	DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
	//评论
	DB.Model(&model.User{}).Select("id,name,created_at,email,avatar,sign,gender").Scan(&users).Count(&total)

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"count": total, "users": users},
		Msg:        "ok",
	}
}

/*********************************************************
** 函数功能: 管理员修改用户信息
** 日    期: 2021年11月12日15:19:08
**********************************************************/
func AdminModifyUserService(newInfo dto.AdminModifyUserRequest) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}
	var user model.User
	DB := common.GetDB()
	//新邮箱的uid不为当前uid
	DB.Where("email = ?", newInfo.Email).First(&user)
	if user.ID != 0 && user.ID != newInfo.ID {
		res.HttpStatus = http.StatusUnprocessableEntity
		res.Code = response.CheckFailCode
		res.Msg = "邮箱已存在"
		return res
	}
	DB.Model(&model.User{}).Where("id = ?", newInfo.ID).Updates(
		map[string]interface{}{
			"email": newInfo.Email,
			"name":  newInfo.Name,
			"sign":  newInfo.Sign,
		},
	)
	return res
}

/*********************************************************
** 函数功能: 管理员删除用户
** 日    期: 2021年11月12日15:26:42
**********************************************************/
func AdminDeleteUserService(id uint) response.ResponseStruct {
	DB := common.GetDB()
	DB.Where("id = ?", id).Delete(model.User{})

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}
}

/*********************************************************
** 函数功能: 邮箱是否被注册
** 日    期: 2021年11月12日11:03:55
**********************************************************/
func IsEmailRegistered(email string) bool {
	DB := common.GetDB()
	return IsEmailExist(DB, email)
}

/*********************************************************
** 函数功能: 邮箱是否属于当前用户
** 日    期: 2021年11月12日11:10:23
**********************************************************/
func IsEmailBelongsToCurrentUser(email string, uid interface{}) bool {
	var user model.User
	DB := common.GetDB()
	DB.First(&user, uid)
	if user.Email == email {
		return true
	}
	return false
}

/*********************************************************
** 函数功能: 邮箱是否存在
** 日    期:2021/7/10
**********************************************************/
func IsEmailExist(db *gorm.DB, email string) bool {
	var user model.User
	db.Where("email = ?", email).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

/*********************************************************
** 函数功能: 用户是否存在
** 日    期:2021/7/10
**********************************************************/
func IsUserExist(db *gorm.DB, id uint) bool {
	var user model.User
	db.First(&user, id)
	if user.ID != 0 {
		return true
	}
	return false
}
