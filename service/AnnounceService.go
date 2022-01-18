package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/common"
	"kuukaa.fun/danmu-v4/dto"
	"kuukaa.fun/danmu-v4/model"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/vo"
)

/*********************************************************
** 函数功能: 获取公告
** 日    期: 2021年11月11日22:06:33
**********************************************************/
func GetAnnounceService(uid interface{}) response.ResponseStruct {
	DB := common.GetDB()
	var oldTime model.AnnounceUser
	var newAnnounce []model.Announce
	//最后一次获取的时间
	DB.Where("uid = ? ", uid.(uint)).Last(&oldTime)
	//查询最后一次获取后的公告
	DB.Where("updated_at > ? ", oldTime.CreatedAt).Find(&newAnnounce)
	if newAnnounce != nil {
		for i := 0; i < len(newAnnounce); i++ {
			DB.Create(&model.AnnounceUser{Aid: newAnnounce[i].ID, Uid: uid.(uint)})
		}
	}
	//拉取用户公告
	var announceList []vo.AnnounceVo
	var sql = "select id,created_at,title,content,url from announces where deleted_at is null and id in (select aid from announce_users where deleted_at is null and uid = ?)"
	DB.Raw(sql, uid).Scan(&announceList)

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"announces": announceList},
		Msg:        "ok",
	}
}

/*********************************************************
** 函数功能: 管理员获取公告
** 日    期: 2021年11月12日12:13:40
**********************************************************/
func AdminGetAnnounceService() response.ResponseStruct {
	DB := common.GetDB()
	var announceList []vo.AdminAnnounceVo
	DB.Raw("select id,created_at,title,content,url from announces where deleted_at is null").Scan(&announceList)
	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"announces": announceList},
		Msg:        "ok",
	}
}

/*********************************************************
** 函数功能: 管理员添加公告
** 日    期: 2021年11月12日12:17:08
**********************************************************/
func AddAnnounceService(announce dto.AddAnnounceDto) response.ResponseStruct {
	newAnnounce := model.Announce{
		Title:   announce.Title,
		Content: announce.Content,
		Url:     announce.Url,
	}
	DB := common.GetDB()
	DB.Create(&newAnnounce)

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}
}

/*********************************************************
** 函数功能: 管理员删除公告
** 日    期: 2021年11月12日12:20:35
**********************************************************/
func DeleteAnnounceService(id uint) response.ResponseStruct {
	DB := common.GetDB()
	DB.Where("id = ?", id).Delete(model.Announce{})

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}
}
