package service

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"kuukaa.fun/danmu-v4/common"
	"kuukaa.fun/danmu-v4/dto"
	"kuukaa.fun/danmu-v4/model"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/util"
	"kuukaa.fun/danmu-v4/vo"
)

/*********************************************************
** 函数功能: 上传视频信息
** 日    期:2021/11/10
** 修改时间: 2021年11月17日12:55:08
** 版    本: 3.5.0
** 修改内容: 移除子视频
** 修改时间: 2021年12月9日17:47:14
** 版    本: 3.6.8
** 修改内容: 分区
**********************************************************/
func UploadVideoInfoService(video dto.UploadVideoDto, uid interface{}) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}

	DB := common.GetDB()
	if !IsSubpartition(DB, video.Partition) {
		res.HttpStatus = http.StatusUnprocessableEntity
		res.Code = response.CheckFailCode
		res.Msg = response.PartitionNotExist
		return res
	}

	newVideo := model.Video{
		Title:       video.Title,
		Cover:       video.Cover,
		Desc:        video.Desc,
		Copyright:   video.Copyright,
		Uid:         uid.(uint),
		VideoType:   viper.GetString("transcoding.coding"),
		PartitionID: video.Partition,
	}

	tx := DB.Begin()
	if err := tx.Create(&newVideo).Error; err != nil {
		tx.Rollback()
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = response.UploadFail
		return res
	}
	//创建审核状态
	if err := tx.Create(&model.Review{Vid: newVideo.ID, Status: 500}).Error; err != nil {
		tx.Rollback()
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = response.UpdateStatusFail
		return res
	}
	tx.Commit()
	res.Data = gin.H{"vid": newVideo.ID}
	return res
}

/*********************************************************
** 函数功能: 获取视频状态
** 日    期: 2021/11/10
**********************************************************/
func GetVideoStatusService(vid int, uid interface{}) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}
	var review model.Review
	DB := common.GetDB()
	DB.Model(&model.Review{}).Preload("Video").Where("vid = ?", vid).First(&review)
	if review.ID == 0 || review.Video.Uid != uid {
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = response.VideoNotExist
		return res
	}
	//通过子分区获取父分区
	partition := GetPartitionName(DB, review.Video.PartitionID)
	var video = vo.ReviewVideoVo{
		Title:     review.Video.Title,
		Cover:     review.Video.Cover,
		Desc:      review.Video.Desc,
		Partition: partition,
	}

	res.Data = gin.H{"status": review.Status, "remarks": review.Remarks, "video": video}
	return res
}

/*********************************************************
** 函数功能: 修改视频信息
** 日    期: 2021/11/10
** 修改时间: 2021年12月9日17:50:35
** 版    本: 3.6.8
** 修改内容: 分区
**********************************************************/
func ModifyVideoInfoService(video dto.ModifyVideoDto, uid interface{}) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}
	DB := common.GetDB()
	tx := DB.Begin()
	if err := tx.Model(&model.Video{}).Where("id = ? and uid = ?", video.ID, uid).Updates(
		map[string]interface{}{
			"title":     video.Title,
			"cover":     video.Cover,
			"desc":      video.Desc,
			"copyright": video.Copyright,
		},
	).Error; err != nil {
		tx.Rollback()
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = response.ModifyFail
		return res
	}
	//更新审核状态
	if err := tx.Model(&model.Review{}).Where("vid = ?", video.ID).Updates(
		map[string]interface{}{
			"status": 1000,
		},
	).Error; err != nil {
		tx.Rollback()
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = response.UpdateStatusFail
		return res
	}
	tx.Commit()
	return res
}

/*********************************************************
** 函数功能: 删除视频
** 日    期:2021/11/10
**********************************************************/
func DeleteVideoService(vid uint, uid interface{}) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}

	DB := common.GetDB()
	if !IsUserOwnsVideo(DB, vid, uid.(uint)) {
		//该视频不属于这个用户
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = response.DeleteFail
		return res
	}
	DB.Where("id = ?", vid).Delete(model.Video{})
	//删除审核状态
	DB.Where("vid = ?", vid).Delete(model.Review{})
	return res
}

/*********************************************************
** 函数功能: 获取自己上传的视频
** 日    期: 2021/11/10
** 修改时间: 2021年11月17日12:59:09
** 版    本: 3.5.0
** 修改内容: 移除子视频
**********************************************************/
func GetMyUploadVideoService(page int, pageSize int, uid interface{}) response.ResponseStruct {
	DB := common.GetDB()
	//记录总数
	var totalSize int
	//分页查询
	var videos []model.Video
	DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
	DB.Where("uid = ?", uid).Find(&videos).Count(&totalSize)
	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"count": totalSize, "data": vo.ToUploadVideoVo(videos)},
		Msg:        response.OK,
	}
}

/*********************************************************
** 函数功能: 通过视频ID获取视频
** 日    期: 2021/11/10
** 修改时间: 2021年11月17日12:59:45
** 版    本: 3.5.0
** 修改内容: 移除子视频
**********************************************************/
func GetVideoByIDService(vid int, ip string, uid interface{}) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}

	var video model.Video
	DB := common.GetDB()
	DB.Model(&model.Video{}).Preload("Author").Where("id = ? and review = true", vid).First(&video)
	if video.ID == 0 {
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.CheckFailCode
		res.Msg = response.VideoNotExist
		return res
	}
	//获取视频资源
	resource := GetVideoResource(DB, uint(vid))
	//获取视频交互数据
	like, collect := CollectAndLikeCount(DB, uint(vid))
	//增加播放量(一个ip在同一个视频下，每30分钟可重新增加1播放量)
	if redis := common.RedisClient; redis != nil {
		clicksLimit, _ := redis.Get(util.VideoClicksLimitKey(vid, ip)).Result()
		if clicksLimit == "" {
			DB.Model(&video).UpdateColumn("clicks", gorm.Expr("clicks + 1"))
			redis.Set(util.VideoClicksLimitKey(vid, ip), 1, time.Minute*30)
		}
	}
	//获取交互数据(如果用户已经登录)
	var interactiveData vo.InteractiveVo
	if uid.(uint) != 0 {
		interactiveData = GetVideoInteractiveData(DB, vid, uid.(uint))
	}

	res.Data = gin.H{
		"video":       vo.ToVideoVo(video, like, collect, resource),
		"interactive": interactiveData,
	}

	return res
}

/*********************************************************
** 函数功能: 获取收藏列表
** 日    期:2021/11/11
**********************************************************/
func GetCollectVideoService(uid interface{}, page int, pageSize int) response.ResponseStruct {
	var count int
	var favorites []model.Interactive

	DB := common.GetDB()
	DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
	DB.Where("uid = ? AND collect = true", uid).Preload("Video").Find(&favorites).Count(&count)

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"count": count, "videos": vo.ToCollectVideoVo(favorites)},
		Msg:        response.OK,
	}
}

/*********************************************************
** 函数功能: 获取推荐视频
** 日    期: 2021/11/11
** 修改时间: 2021年11月17日13:00:01
** 版    本: 3.5.0
** 修改内容: 移除子视频
**********************************************************/
func GetRecommendVideoService() response.ResponseStruct {
	var videos []vo.RecommendVideoVo
	DB := common.GetDB()
	DB = DB.Limit(8)

	const sql = "select videos.id,title,cover,name as author,clicks from" +
		" users,videos where users.id=videos.uid and review=1 and videos.deleted_at is null order by clicks desc"

	DB.Raw(sql).Scan(&videos)

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"videos": videos},
		Msg:        response.OK,
	}
}

/*********************************************************
** 函数功能: 获取视频列表
** 日    期: 2021/11/11
** 修改时间: 2021年11月17日13:00:26
** 版    本: 3.5.0
** 修改内容: 移除子视频
**********************************************************/
func GetVideoListService(query dto.GetVideoListDto) response.ResponseStruct {
	DB := common.GetDB()
	var total int //记录总数
	var videos []vo.SearchVideoVo
	Pagination := DB.Limit(query.PageSize).Offset((query.Page - 1) * query.PageSize)

	if query.Partition == 0 {
		//不传分区参数默认查询全部
		Pagination.Model(&model.Video{}).Select("id,title,cover").Where("review = 1").Scan(&videos).Count(&total)
	} else if IsSubpartition(DB, uint(query.Partition)) {
		//判断是否为子分区
		Pagination.Model(&model.Video{}).Select("id,title,cover").Where(
			map[string]interface{}{"review": 1, "partition_id": query.Partition},
		).Scan(&videos).Count(&total)
	} else {
		//获取该分区下的子分区
		list := GetSubpartitionList(DB, uint(query.Partition))
		Pagination.Debug().Model(&model.Video{}).Select("id,title,cover").
			Where("review = 1 and partition_id in (?)", list).Scan(&videos).Count(&total)
	}

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"count": total, "videos": videos},
		Msg:        response.OK,
	}
}

/*********************************************************
** 函数功能: 通过用户ID获取视频列表
** 日    期: 2021/11/11
** 修改时间: 2021年11月17日13:00:41
** 版    本: 3.5.0
** 修改内容: 移除子视频
**********************************************************/
func GetVideoListByUserIDService(uid int, page int, pageSize int) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}
	var videos []vo.SearchVideoVo

	DB := common.GetDB()
	if !IsUserExist(DB, uint(uid)) {
		res.HttpStatus = http.StatusUnprocessableEntity
		res.Code = response.CheckFailCode
		res.Msg = response.UserNotExist
		return res
	}
	//记录总数
	var total int
	DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
	DB.Model(&model.Video{}).Select("id,title,cover").Where("review = 1 and uid = ?", uid).Scan(&videos).Count(&total)
	res.Data = gin.H{"count": total, "videos": videos}
	return res
}

/*********************************************************
** 函数功能: 管理员获取视频列表
** 日    期: 2021年11月12日15:30:26
**********************************************************/
func AdminGetVideoListService(page int, pageSize int, videoFrom string) response.ResponseStruct {
	var total int //记录总数
	var videos []model.Video

	DB := common.GetDB()
	DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
	if videoFrom == "admin" {
		DB.Where("review = 1 and uid = 0").Find(&videos).Count(&total)
	} else {
		DB.Where("review = 1 and uid != 0").Find(&videos).Count(&total)
	}

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"count": total, "videos": vo.ToAdminVideoListVo(videos)},
		Msg:        response.OK,
	}
}

/*********************************************************
** 函数功能: 管理员删除视频
** 日    期: 2021年11月12日15:33:20
**********************************************************/
func AdminDeleteVideoService(id uint) response.ResponseStruct {
	DB := common.GetDB()
	DB.Where("id = ?", id).Delete(model.Video{})
	//删除审核状态
	DB.Where("vid = ?", id).Delete(model.Review{})
	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}
}

/*********************************************************
** 函数功能: 管理员导入视频
** 日    期: 2021年11月12日15:36:28
**********************************************************/
func ImportVideoService(video dto.ImportVideo) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}

	newVideo := model.Video{
		Title:     video.Title,
		Cover:     video.Cover,
		Desc:      video.Desc,
		Copyright: true,
		Uid:       uint(viper.GetInt("user.video")),
		VideoType: video.Type,
		Review:    true,
	}

	DB := common.GetDB()
	if err := DB.Create(&newVideo).Error; err != nil {
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.CheckFailCode
		res.Msg = response.UploadFail
		return res
	}

	return res
}

/*********************************************************
** 函数功能: 管理员导入视频资源
** 日    期: 2022年1月13日16:22:05
**********************************************************/
func ImportResourceService(video dto.ImportResourceDto) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}

	DB := common.GetDB()
	//添加视频链接
	if err := DB.Model(&model.Resource{}).Create(&model.Resource{
		Vid:      video.Vid,
		Res360:   video.Res360,
		Res480:   video.Res480,
		Res720:   video.Res720,
		Res1080:  video.Res1080,
		Original: video.Original,
	}).Error; err != nil {
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = response.CreateFail
		return res
	}

	return res
}

/*********************************************************
** 函数功能: 管理员获取视频资源
** 日    期: 2022年1月14日11:29:45
**********************************************************/
func GetResourceListService(vid int) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}

	DB := common.GetDB()
	var resources []vo.ResourceInfoVo
	DB.Model(&model.Resource{}).Select("uuid,title,created_at").Where("vid = ?", vid).Scan(&resources)

	res.Data = gin.H{"resources": resources}
	return res
}

/*********************************************************
** 函数功能: 管理员删除视频资源
** 日    期: 2022年1月14日12:11:37
**********************************************************/
func DeleteResourceService(uuid uuid.UUID) response.ResponseStruct {
	DB := common.GetDB()
	DB.Where("uuid = ?", uuid).Delete(model.Resource{})
	//删除审核状态
	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}
}

/*********************************************************
** 函数功能: 视频是否属于自己
** 日    期:2021/11/6
**********************************************************/
func IsUserOwnsVideo(db *gorm.DB, vid uint, uid uint) bool {
	var video model.Video
	db.Where("id = ? and uid = ?", vid, uid).First(&video)
	if video.ID != 0 {
		return true
	}
	return false
}

/*********************************************************
** 函数功能: 视频是否存在
** 日    期:2021/7/22
**********************************************************/
func IsVideoExist(db *gorm.DB, vid uint) bool {
	var video model.Video
	db.First(&video, vid)
	if video.ID != 0 {
		return true
	}
	return false
}

/*********************************************************
** 函数功能: 获取视频资源
** 日    期: 2022年1月6日10:33:53
**********************************************************/
func GetVideoResource(db *gorm.DB, vid uint) []model.Resource {
	var resource []model.Resource
	db.Model(&model.Resource{}).Where("vid = ?", vid).Find(&resource)
	return resource
}

/*********************************************************
** 函数功能: 获取视频交互数据
** 日    期:2021/11/11
**********************************************************/
func GetVideoInteractiveData(db *gorm.DB, vid int, uid uint) vo.InteractiveVo {
	var video model.Video
	//获取作者id
	db.First(&video, vid)

	like, collect := IsCollectAndLike(db, uid, uint(vid))
	follow := IsFollow(db, uid, video.Uid)

	return vo.InteractiveVo{
		Collect: collect,
		Like:    like,
		Follow:  follow,
	}
}
