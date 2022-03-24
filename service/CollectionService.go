package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"kuukaa.fun/danmu-v4/common"
	"kuukaa.fun/danmu-v4/dto"
	"kuukaa.fun/danmu-v4/model"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/vo"
)

/*********************************************************
** 函数功能: 创建视频合集
** 日    期: 2021/11/19
**********************************************************/
func CreateCollectionService(collection dto.CreateCollectionDto, uid interface{}) response.ResponseStruct {
	DB := common.GetDB()

	newCollection := model.Collection{
		Title: collection.Title,
		Cover: collection.Cover,
		Desc:  collection.Desc,
		Uid:   uid.(uint),
	}

	DB.Create(&newCollection)
	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}
}

/*********************************************************
** 函数功能: 修改合集信息
** 日    期: 2021/11/21
**********************************************************/
func ModifyCollectionService(collection dto.ModifyCollectionDto, uid interface{}) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}

	DB := common.GetDB()

	if !isUserOwnsCollection(DB, collection.ID, uid.(uint)) {
		res.HttpStatus = http.StatusUnprocessableEntity
		res.Code = response.CheckFailCode
		res.Msg = response.CollectionNotExist
		return res
	}

	err := DB.Model(model.Collection{}).Where("id = ?", collection.ID).Updates(
		map[string]interface{}{"cover": collection.Cover, "title": collection.Title, "desc": collection.Desc},
	).Error
	if err != nil {
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = response.ModifyFail
		return res
	}
	return res
}

/*********************************************************
** 函数功能: 获取合集列表
** 日    期: 2021/11/19
**********************************************************/
func GetCreateCollectionListService(page int, pageSize int, uid interface{}) response.ResponseStruct {
	var count int
	var collections []vo.CollectionVo
	DB := common.GetDB()
	DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
	DB.Raw("select id,title,cover,`desc`,created_at from collections where deleted_at is null and uid = ?", uid).Scan(&collections)
	DB.Model(&model.Collection{}).Where("uid = ?", uid).Count(&count)
	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"count": count, "collections": collections},
		Msg:        response.OK,
	}
}

/*********************************************************
** 函数功能: 获取合集内容
** 日    期: 2021/11/19
**********************************************************/
func GetCollectionContentService(cid int, page int, pageSize int) response.ResponseStruct {
	var count int
	var videos []vo.CollectionVideoVo
	DB := common.GetDB()
	DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
	sqlID := "(select vid from video_collections where deleted_at is null and collection_id = ?)"
	DB.Raw("select id,title,cover,created_at,`desc` from videos where deleted_at is null and review = 1 and id in "+sqlID, cid).Scan(&videos)
	DB.Model(&model.VideoCollection{}).Where("collection_id = ?", cid).Count(&count) //获取已添加数量

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"count": count, "videos": videos},
		Msg:        response.OK,
	}
}

/*********************************************************
** 函数功能: 获取合集信息
** 日    期: 2021/11/20
**********************************************************/
func GetCollectionByIDService(id int) response.ResponseStruct {
	var user model.User
	var collection model.Collection
	DB := common.GetDB()

	DB.Model(&model.Collection{}).Where("id = ?", id).First(&collection)
	DB.Model(&model.User{}).Where("id = ?", collection.Uid).First(&user)

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"collection": vo.ToCollectionVo(collection), "user": vo.ToAuthorVo(user)},
		Msg:        response.OK,
	}
}

/*********************************************************
** 函数功能: 删除合集
** 日    期:2021/11/20
**********************************************************/
func DeleteCollectionService(id uint, uid interface{}) response.ResponseStruct {
	DB := common.GetDB()
	DB.Where("id = ? and uid = ?", id, uid).Delete(model.Collection{})
	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}
}

/*********************************************************
** 函数功能: 获取可添加视频
** 日    期: 2021/11/20
**********************************************************/
func GetCanAddVideoService(id int, uid interface{}, page int, pageSize int) response.ResponseStruct {
	var count int
	var addedCount int
	var videos []vo.CollectionVideoVo

	DB := common.GetDB()
	DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
	sql := "select id,title,cover from videos where deleted_at is null and uid = ? and review = 1 and id not in "
	sqlVid := "(select vid from video_collections where deleted_at is null and collection_id = ?)"

	DB.Raw(sql+sqlVid, uid, id).Scan(&videos)                                            //查询可添加视频列表
	DB.Model(&model.Video{}).Where("review = 1 and uid = ?", uid).Count(&count)          //获取视频总数
	DB.Model(&model.VideoCollection{}).Where("collection_id = ?", id).Count(&addedCount) //获取已添加数量

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"count": count - addedCount, "videos": videos},
		Msg:        response.OK,
	}
}

/*********************************************************
** 函数功能: 向合集内添加视频
** 日    期: 2021/11/21
**********************************************************/
func AddVideoToCollectionService(vid uint, cid uint, uid interface{}) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}

	DB := common.GetDB()
	if !isUserOwnsVideo(DB, vid, uid.(uint)) {
		res.HttpStatus = http.StatusUnprocessableEntity
		res.Code = response.CheckFailCode
		res.Msg = response.VideoNotExist
		return res
	}

	if !isUserOwnsCollection(DB, cid, uid.(uint)) {
		res.HttpStatus = http.StatusUnprocessableEntity
		res.Code = response.CheckFailCode
		res.Msg = response.PartitionNotExist
		return res
	}

	//添加视频
	newVideoCollection := model.VideoCollection{
		Vid:          vid,
		CollectionId: cid,
	}

	DB.Create(&newVideoCollection)
	return res
}

/*********************************************************
** 函数功能: 删除合集内的视频
** 日    期: 2021年11月21日13:33:15
**********************************************************/
func DeleteCollectionVideoService(vid uint, cid uint, uid interface{}) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}

	DB := common.GetDB()
	if !isUserOwnsVideo(DB, vid, uid.(uint)) {
		res.HttpStatus = http.StatusUnprocessableEntity
		res.Code = response.CheckFailCode
		res.Msg = response.VideoNotExist
		return res
	}

	if !isUserOwnsCollection(DB, cid, uid.(uint)) {
		res.HttpStatus = http.StatusUnprocessableEntity
		res.Code = response.CheckFailCode
		res.Msg = response.PartitionNotExist
		return res
	}

	if err := DB.Where("collection_id = ? and vid = ?", cid, vid).Delete(model.VideoCollection{}).Error; err != nil {
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = response.DeleteFail
		return res
	}
	return res
}

/*********************************************************
** 函数功能: 获取合集列表
** 日    期: 2021/11/21
**********************************************************/
func GetCollectionListService(page int, pageSize int) response.ResponseStruct {
	var count int
	var collections []vo.CollectionVo
	DB := common.GetDB()
	DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
	DB.Raw("select id,title,cover,`desc`,created_at from collections where deleted_at is null").Scan(&collections)
	DB.Model(&model.Collection{}).Count(&count)
	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"count": count, "collections": collections},
		Msg:        response.OK,
	}
}

/*********************************************************
** 函数功能: 管理员删除合集
** 日    期: 2022年2月24日15:16:37
**********************************************************/
func AdminDeleteCollectionService(id uint) response.ResponseStruct {
	DB := common.GetDB()
	DB.Where("id = ?", id).Delete(model.Collection{})

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}
}

/*********************************************************
** 函数功能: 合集是否属于用户
** 日    期: 2021/11/21
**********************************************************/
func isUserOwnsCollection(db *gorm.DB, cid uint, uid uint) bool {
	var collection model.Collection
	db.Where("id = ? and uid = ?", cid, uid).First(&collection)
	if collection.ID != 0 {
		return true
	}
	return false
}
