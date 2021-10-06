package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/util"
)

/*********************************************************
** 函数功能: 用户注册
** 日    期:2021/7/10
**********************************************************/
func Register(ctx *gin.Context) {
	//获取参数
	type requestRegister struct {
		Name string
		Email string
		Password string
		Code string
	}
	var requestUser requestRegister
	err := ctx.Bind(&requestUser)
	if err != nil{
		response.Fail(ctx,nil,"请求错误")
		return
	}
	name := requestUser.Name
	email := requestUser.Email
	password := requestUser.Password
	code := requestUser.Code
	//数据验证
	if !util.VerifyEmailFormat(email) {
		response.CheckFail(ctx,nil,"邮箱格式有误哦")
		return
	}
	if len(password) < 6 {
		response.CheckFail(ctx,nil,"密码不要少于六位")
		return
	}
	if !VerificationCode(email,code){
		response.CheckFail(ctx,nil,"验证码有误")
		return
	}
	//邮箱是否存在
	DB :=common.GetDB()
	if IsEmailExist(DB,email){
		response.CheckFail(ctx,nil,"该邮箱已经被注册了")
		return
	}
	//如果名称为空，则为随机字符串
	if len(name) == 0{
		name = util.RandomString(10)
	}

	//创建用户
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil {
		response.CheckFail(ctx,nil,"服务器出错了")
		//记录日志
		util.Logfile("[Error]"," hashed password " + err.Error())
		return
	}

	newUser := model.User{
		Name:name,
		Email:email,
		Password:string(hashedPassword),
	}
	DB.Create(&newUser)
	//返回结果
	response.Success(ctx,nil,"注册成功")
}


/*********************************************************
** 函数功能: 用户登录
** 日    期:2021/7/10
**********************************************************/
func Login(ctx *gin.Context)  {
	//获取参数
	var user = model.User{}
	requestErr := ctx.Bind(&user)
	if requestErr != nil{
		response.Response(ctx,http.StatusBadRequest,4000,nil,"请求错误")
		return
	}
	password := user.Password
	email := user.Email
	//数据验证
	if !util.VerifyEmailFormat(email) {
		response.CheckFail(ctx,nil,"邮箱格式有误哦")
		return
	}
	if len(password) < 6 {
		response.CheckFail(ctx,nil,"密码不要少于六位")
		return
	}
	//判断邮箱是否存在
	DB :=common.GetDB()
	DB.Where("email = ?" ,email).First(&user)
	if user.ID == 0{
		response.Fail(ctx,nil,"用户不存在")
		return
	}
	//判断密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(password)) ; err != nil {
		response.CheckFail(ctx,nil,"用户名或密码错误")
		return
	}
	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.ServerError(ctx,nil,"系统异常")
		util.Logfile("[Error]"," token generate error  " + err.Error())
		return
	}
	//返回数据
	response.Success(ctx,gin.H{"token":token,"user":dto.ToUserDto(user)},"登陆成功")
}


/*********************************************************
** 函数功能: 用户获取个人信息
** 日    期:2021/7/10
**********************************************************/
func UserInfo(ctx *gin.Context)  {
	user,_ :=ctx.Get("user")
	response.Success(ctx,gin.H{"data":dto.ToUserDto(user.(model.User))},"")
}

/*********************************************************
** 函数功能: 用户修改个人信息
** 日    期:2021/7/10
**********************************************************/
func ModifyInfo(ctx *gin.Context)  {
	//获取参数
	type modifyRegister struct {
		Name string
		Gender int
		Birthday string
		Sign string
	}
	var user = modifyRegister{}
	err := ctx.Bind(&user)
	if err != nil{
		response.Fail(ctx,nil,"请求错误")
		return
	}
	name := user.Name
	gender := user.Gender
	birthday := user.Birthday
	sign := user.Sign
	if len(name)==0{
		response.CheckFail(ctx,nil,"昵称不能为空哦")
		return
	}
	//判断日期
	tBirthday,err := time.Parse("2006-01-02",birthday)
	if err != nil{
		response.CheckFail(ctx,nil,"请输入正确的出生日期哦")
		return
	}

	//从上下文中获取用户id
	id,_ :=ctx.Get("id")
	DB :=common.GetDB()
	err = DB.Model(model.User{}).Where("id = ?",id).Updates(map[string]interface{}{"name":name,"gender":gender,"birthday":tBirthday,"sign":sign}).Error
	if err != nil{
		response.Fail(ctx,nil,"修改失败")
	}else {
		response.Success(ctx,nil,"ok")
	}
}

/*********************************************************
** 函数功能: 通过用户ID获取用户信息
** 日    期:2021/7/10
**********************************************************/
func GetUserInfoByID(ctx *gin.Context)  {
	var user model.User
	uid, _ := strconv.Atoi(ctx.Query("uid"))
	DB :=common.GetDB()
	DB.Select("id,name,sign,avatar,gender").Where("id = ?",uid).First(&user)
	response.Success(ctx,gin.H{"user":dto.ToUserDto(user)},"ok")
}

/*********************************************************
** 函数功能: 邮箱是否存在
** 日    期:2021/7/10
**********************************************************/
func IsEmailExist(db *gorm.DB,email string) bool {
	var user model.User
	db.Where("email = ?" ,email).First(&user)
	if user.ID != 0{
		return true
	}
	return false
}

/*********************************************************
** 函数功能: 用户是否存在
** 日    期:2021/7/10
**********************************************************/
func IsUserExist(db *gorm.DB,id uint) bool {
	var user model.User
	db.First(&user,id)
	if user.ID != 0{
		return true
	}
	return false
}