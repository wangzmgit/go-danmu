package controller

import (
	"os"
	"path"
	"strconv"
	"time"
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

/**********************************************用户相关接口*************************************************/

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
		DB.Model(&model.User{}).Select("id,name,email,avatar,sign,gender").Scan(&users).Count(&total)
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

/**********************************************视频相关接口*************************************************/

/*********************************************************
** 函数功能: 获取视频列表
** 日    期:2021/8/4
**********************************************************/
func AdminGetVideoList(ctx *gin.Context) {
	DB := common.GetDB()
	var videos []model.Video
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	if page > 0 && pageSize > 0 {
		//记录总数
		var total int
		DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
		DB.Where("review = 1").Find(&videos).Count(&total)
		response.Success(ctx, gin.H{"count": total, "videos": dto.ToAdminVideoDto(videos)}, "ok")
	} else {
		response.Fail(ctx, nil, "获取数量有误")
	}
}

/*********************************************************
** 函数功能: 删除视频
** 日    期:2021/8/3
**********************************************************/
func AdminDeleteVideo(ctx *gin.Context) {
	DB := common.GetDB()
	var request = AdminIDRequest{}
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	id := request.ID
	DB.Where("id = ?", id).Delete(model.Video{})
	response.Success(ctx, nil, "ok")
}

/**********************************************审核相关接口*************************************************/
/*********************************************************
** 函数功能: 获取待审核视频列表
** 日    期:2021/8/4
**********************************************************/
func GetReviewVideoList(ctx *gin.Context) {
	DB := common.GetDB()
	var videos []model.Video
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	if page > 0 && pageSize > 0 {
		//记录总数
		var total int
		DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
		DB.Raw("select * from videos where deleted_at is null and id in (select vid from reviews where deleted_at is null and status = 1000)").Scan(&videos)
		response.Success(ctx, gin.H{"count": total, "videos": dto.ToAdminVideoDto(videos)}, "ok")
	} else {
		response.Fail(ctx, nil, "获取数量有误")
	}
}

/*********************************************************
** 函数功能: 审核视频
** 日    期:2021/8/4
**********************************************************/
func ReviewVideo(ctx *gin.Context) {
	type review struct {
		VID     uint
		Status  int
		Remarks string
	}
	var requestReview review
	err := ctx.Bind(&requestReview)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	vid := requestReview.VID
	status := requestReview.Status
	remarks := requestReview.Remarks
	var isReview bool
	if vid == 0 {
		response.CheckFail(ctx, nil, "视频不存在")
		return
	}
	if status == 2000 {
		isReview = true
	} else if status == 4001 || status == 4002 {
		isReview = false
	} else {
		response.CheckFail(ctx, nil, "状态错误")
		return
	}
	DB := common.GetDB()
	tx := DB.Begin()
	if err := tx.Model(&model.Video{}).Where("id = ?", vid).Updates(map[string]interface{}{"review": isReview}).Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "修改失败")
		return
	}
	//创建审核状态
	if err := tx.Model(&model.Review{}).Where("vid = ?", vid).Updates(map[string]interface{}{"status": status, "remarks": remarks}).Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "状态更新失败")
		return
	}
	tx.Commit()
	response.Success(ctx, nil, "ok")
}

/**********************************************公告相关接口*************************************************/

/*********************************************************
** 函数功能: 获取公告
** 日    期:2021/8/4
**********************************************************/
func AdminGetAnnounce(ctx *gin.Context) {
	DB := common.GetDB()
	var announceList []dto.AdminAnnounceDto
	DB.Raw("select id,created_at,title,content,url from announces where deleted_at is null").Scan(&announceList)
	response.Success(ctx, gin.H{"announces": announceList}, "ok")
}

/*********************************************************
** 函数功能: 添加公告
** 日    期:2021/8/4
**********************************************************/
func AddAnnounce(ctx *gin.Context) {
	DB := common.GetDB()
	var announce = model.Announce{}
	err := ctx.Bind(&announce)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	title := announce.Title
	content := announce.Content
	url := announce.Url
	if len(title) == 0 {
		response.CheckFail(ctx, nil, "标题不能为空")
		return
	}
	if len(content) == 0 {
		response.CheckFail(ctx, nil, "内容不能为空")
		return
	}
	newAnnounce := model.Announce{
		Title:   title,
		Content: content,
		Url:     url,
	}
	DB.Create(&newAnnounce)
	//返回结果
	response.Success(ctx, nil, "ok")
}

/*********************************************************
** 函数功能: 删除公告
** 日    期:2021/8/4
**********************************************************/
func DeleteAnnounce(ctx *gin.Context) {
	DB := common.GetDB()
	var request = AdminIDRequest{}
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	id := request.ID
	DB.Where("id = ?", id).Delete(model.Announce{})
	response.Success(ctx, nil, "ok")
}

/**********************************************轮播图相关接口*************************************************/

/*********************************************************
** 函数功能: 上传轮播图
** 日    期:2021/8/4
**********************************************************/
func UploadCarousel(ctx *gin.Context) {
	carousel, err := ctx.FormFile("carousel")
	if err != nil {
		response.Fail(ctx, nil, "图片上传失败")
		return
	}
	suffix := path.Ext(carousel.Filename)
	if suffix != ".jpg" && suffix != ".jpeg" && suffix != ".png" {
		response.CheckFail(ctx, nil, "图片不符合要求")
		return
	}
	carousel.Filename = util.RandomString(3) + strconv.FormatInt(time.Now().UnixNano(), 10) + suffix
	errSave := ctx.SaveUploadedFile(carousel, "./file/carousel/"+carousel.Filename)
	if errSave != nil {
		response.Fail(ctx, nil, "图片保存失败")
		return
	}
	fileInfo, err := os.Stat("./file/carousel/" + carousel.Filename)
	//大小限制到5M
	if fileInfo == nil || fileInfo.Size() > 1024*1024*5 || err != nil {
		response.CheckFail(ctx, nil, "图片大小不符合要求")
		return
	}
	// 拼接上传图片的路径信息
	localFileName := "./file/carousel/" + carousel.Filename
	objectName := "carousel/" + carousel.Filename
	success, url := util.UploadOSS(localFileName, objectName)
	if success {
		response.Success(ctx, gin.H{"url": url}, "ok")
	} else {
		response.Fail(ctx, nil, "上传失败")
	}
}

/*********************************************************
** 函数功能: 上传轮播图信息
** 日    期:2021/8/4
**********************************************************/
func UploadCarouselInfo(ctx *gin.Context) {
	DB := common.GetDB()
	var carousel model.Carousel
	err := ctx.Bind(&carousel)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	img := carousel.Img
	url := carousel.Url
	//验证数据
	if len(img) == 0 {
		response.CheckFail(ctx, nil, "图片不能为空")
		return
	}
	newCarousel := model.Carousel{
		Img: img,
		Url: url,
	}
	DB.Create(&newCarousel)
	response.Success(ctx, nil, "ok")
}

/*********************************************************
** 函数功能: 获取轮播图
** 日    期:2021/8/4
**********************************************************/
func AdminGetCarousel(ctx *gin.Context) {
	DB := common.GetDB()
	var carousels []dto.AdminCarouselDto
	DB.Model(&model.Carousel{}).Select("id,img,url,created_at").Scan(&carousels)
	response.Success(ctx, gin.H{"carousels": carousels}, "ok")
}

/*********************************************************
** 函数功能: 删除轮播图
** 日    期:2021/8/4
**********************************************************/
func DeleteCarousel(ctx *gin.Context) {
	DB := common.GetDB()
	var request = AdminIDRequest{}
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	id := request.ID
	DB.Where("id = ?", id).Delete(model.Carousel{})
	response.Success(ctx, nil, "ok")
}
