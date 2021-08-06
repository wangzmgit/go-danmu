package controller

import (
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"

	"github.com/gin-gonic/gin"
)

/*********************************************************
** 函数功能: 获取公告
** 日    期:2021/7/29
**********************************************************/
func GetAnnounce(ctx *gin.Context) {
	DB := common.GetDB()
	var oldTime model.AnnounceUser
	var newAnnounce []model.Announce

	uid, _ := ctx.Get("id")
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
	var announceList []dto.AnnounceDto
	var sql = "select id,created_at,title,content,url from announces where deleted_at is null and id in (select aid from announce_users where deleted_at is null and uid = ?)"
	DB.Raw(sql, uid).Scan(&announceList)
	response.Success(ctx, gin.H{"announces": announceList}, "ok")
}
