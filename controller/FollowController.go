package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"
)

/*********************************************************
** 函数功能: 关注
** 日    期:2021/7/24
**********************************************************/
func Following(ctx *gin.Context){
	//获取参数
	var follow = model.Follow{}
	err := ctx.Bind(&follow)
	if err != nil{
		response.Response(ctx,http.StatusBadRequest,400,nil,"请求错误")
		return
	}
	//关注的人的id和自己的id
	fid := follow.ID
	uid,_ :=ctx.Get("id")
	//判断关注的是否为自己
	if fid == uid{
		response.CheckFail(ctx,nil,"不能关注自己")
		return
	}
	//验证关注的人是否存在
	DB :=common.GetDB()
	if !IsUserExist(DB,fid){
		response.CheckFail(ctx,nil,"用户不存在")
		return
	}
	//没有记录自动创建记录
	var followInfo model.Follow
	DB.Where("uid = ? and fid = ?",uid,fid).Attrs(model.Follow{Uid:uid.(uint), Fid:fid,}).FirstOrCreate(&followInfo)
	response.Success(ctx,nil,"ok")
}

/*********************************************************
** 函数功能: 取消关注
** 日    期:2021/7/24
**********************************************************/
func UnFollow(ctx *gin.Context)  {
	var follow = model.Follow{}
	err := ctx.Bind(&follow)
	if err != nil{
		response.Response(ctx,http.StatusBadRequest,400,nil,"请求错误")
		return
	}
	//关注的人的id和自己的id
	fid := follow.ID
	uid,_ :=ctx.Get("id")
	DB :=common.GetDB()
	DB.Where("uid = ? and fid = ?",uid,fid).Delete(model.Follow{})
	response.Success(ctx,nil,"ok")
}

/*********************************************************
** 函数功能: 获取关注状态
** 日    期:2021/7/25
**********************************************************/
func GetFollowStatus(ctx *gin.Context)  {
	fid, _ := strconv.Atoi(ctx.Query("fid"))
	if fid == 0 {
		response.CheckFail(ctx, nil, "用户不存在")
		return
	}
	DB :=common.GetDB()
	uid,_ :=ctx.Get("id")
	follow := IsFollow(DB,uid.(uint),uint(fid))
	response.Success(ctx,gin.H{"follow":follow},"ok")
}

/*********************************************************
** 函数功能: 通过UID获取关注列表
** 日    期:2021/7/25
**********************************************************/
func GetFollowingByID(ctx *gin.Context)  {
	uid, _ := strconv.Atoi(ctx.Query("uid"))
	if uid == 0 {
		response.CheckFail(ctx, nil, "用户不存在")
		return
	}
	var users []dto.FollowDto
	DB :=common.GetDB()
	DB.Raw("select id,name,sign,avatar from users where deleted_at is null and id in (select fid from follows where uid = ? and deleted_at is null)",uid).Scan(&users)
	response.Success(ctx,gin.H{"users":users},"ok")
}

/*********************************************************
** 函数功能: 通过UID获取粉丝列表
** 日    期:2021/7/25
**********************************************************/
func GetFollowersByID(ctx *gin.Context)  {
	uid, _ := strconv.Atoi(ctx.Query("uid"))
	if uid == 0 {
		response.CheckFail(ctx, nil, "用户不存在")
		return
	}
	var users []dto.FollowDto
	DB :=common.GetDB()
	DB.Raw("select id,name,sign,avatar from users where deleted_at is null and id in (select uid from follows where fid = ? and deleted_at is null)",uid).Scan(&users)
	response.Success(ctx,gin.H{"users":users},"ok")
}

/*********************************************************
** 函数功能: 通过UID获取粉丝数
** 日    期:2021/7/25
**********************************************************/
func GetFollowCount(ctx *gin.Context)  {
	uid, _ := strconv.Atoi(ctx.Query("uid"))
	if uid == 0 {
		response.CheckFail(ctx, nil, "用户不存在")
		return
	}
	var following int
	var followers int
	DB :=common.GetDB()
	DB.Model(&model.Follow{}).Where("uid = ?",uid).Count(&following)
	DB.Model(&model.Follow{}).Where("fid = ?",uid).Count(&followers)
	response.Success(ctx,gin.H{"following":following,"followers":followers},"ok")
}

/*********************************************************
** 函数功能: 是否关注
** 日    期:2021/7/24
**********************************************************/
func IsFollow(db *gorm.DB,uid uint,fid uint) bool {
	var follow model.Follow
	db.Where("uid = ? and fid = ?", uid,fid).First(&follow)
	if follow.ID != 0 {
		return true
	}
	return false
}
