package service

import (
	"net/http"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/vo"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

/*********************************************************
** 函数功能: 关注
** 日    期:2021/11/10
**********************************************************/
func FollowingService(fid uint, uid interface{}) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}
	//验证关注的人是否存在
	DB := common.GetDB()
	if !IsUserExist(DB, fid) {
		res.HttpStatus = http.StatusUnprocessableEntity
		res.Code = response.CheckFailCode
		res.Msg = "用户不存在"
		return res
	}
	//没有记录自动创建记录
	var followInfo model.Follow
	DB.Where("uid = ? and fid = ?", uid, fid).Attrs(model.Follow{Uid: uid.(uint), Fid: fid}).FirstOrCreate(&followInfo)
	return res
}

/*********************************************************
** 函数功能: 取消关注
** 日    期:2021/11/11
**********************************************************/
func UnFollowService(fid uint, uid interface{}) response.ResponseStruct {
	DB := common.GetDB()
	DB.Where("uid = ? and fid = ?", uid, fid).Delete(model.Follow{})
	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}
}

/*********************************************************
** 函数功能: 获取关注状态
** 日    期:2021/11/11
**********************************************************/
func GetFollowStatusService(uid uint, fid uint) response.ResponseStruct {
	DB := common.GetDB()
	follow := IsFollow(DB, uid, fid)
	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"follow": follow},
		Msg:        "ok",
	}
}

/*********************************************************
** 函数功能: 通过UID获取关注列表
** 日    期:2021/11/11
**********************************************************/
func GetFollowingByIDService(uid interface{}, page int, pageSize int) response.ResponseStruct {
	var users []vo.FollowVo
	DB := common.GetDB()
	DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
	DB.Raw("select id,name,sign,avatar from users where deleted_at is null and id in (select fid from follows where uid = ? and deleted_at is null)", uid).Scan(&users)
	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"users": users},
		Msg:        "ok",
	}
}

/*********************************************************
** 函数功能: 通过UID获取粉丝列表
** 日    期:2021/11/11
**********************************************************/
func GetFollowersByIDService(uid interface{}, page int, pageSize int) response.ResponseStruct {
	var users []vo.FollowVo
	DB := common.GetDB()
	DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
	DB.Raw("select id,name,sign,avatar from users where deleted_at is null and id in (select uid from follows where fid = ? and deleted_at is null)", uid).Scan(&users)
	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"users": users},
		Msg:        "ok",
	}
}

/*********************************************************
** 函数功能: 通过UID获取粉丝数
** 日    期:2021/7/25
**********************************************************/
func GetFollowCountService(uid int) response.ResponseStruct {
	var following int
	var followers int
	DB := common.GetDB()
	DB.Model(&model.Follow{}).Where("uid = ?", uid).Count(&following)
	DB.Model(&model.Follow{}).Where("fid = ?", uid).Count(&followers)
	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"following": following, "followers": followers},
		Msg:        "ok",
	}
}

/*********************************************************
** 函数功能: 是否关注
** 日    期:2021/7/24
**********************************************************/
func IsFollow(db *gorm.DB, uid uint, fid uint) bool {
	var follow model.Follow
	db.Where("uid = ? and fid = ?", uid, fid).First(&follow)
	if follow.ID != 0 {
		return true
	}
	return false
}
