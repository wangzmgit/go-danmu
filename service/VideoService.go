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
	"github.com/spf13/viper"
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
func UploadVideoInfoService(video dto.UploadVideoRequest, uid interface{}) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}

	DB := common.GetDB()
	if !IsSubpartition(DB, video.Partition) {
		res.HttpStatus = http.StatusUnprocessableEntity
		res.Code = response.CheckFailCode
		res.Msg = "分区不存在"
		return res
	}
	newVideo := model.Video{
		Title:        video.Title,
		Cover:        video.Cover,
		Introduction: video.Introduction,
		Original:     video.Original,
		Uid:          uid.(uint),
		VideoType:    viper.GetString("server.coding"),
		PartitionID:  video.Partition,
	}

	tx := DB.Begin()
	if err := tx.Create(&newVideo).Error; err != nil {
		tx.Rollback()
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = "上传失败"
		return res
	}
	//创建审核状态
	if err := tx.Create(&model.Review{Vid: newVideo.ID, Status: 500}).Error; err != nil {
		tx.Rollback()
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = "更新审核状态失败"
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
		Msg:        "ok",
	}
	var review model.Review
	DB := common.GetDB()
	DB.Model(&model.Review{}).Preload("Video").Where("vid = ?", vid).First(&review)
	if review.ID == 0 || review.Video.Uid != uid {
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = "视频不见了"
		return res
	}
	//通过子分区获取父分区
	partition := GetPartitionName(DB, review.Video.PartitionID)
	var video = vo.ReviewVideoVo{
		Title:        review.Video.Title,
		Cover:        review.Video.Cover,
		Introduction: review.Video.Introduction,
		Partition:    partition,
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
func ModifyVideoInfoService(video dto.VideoModifyRequest, uid interface{}) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}
	DB := common.GetDB()
	tx := DB.Begin()
	if err := tx.Model(&model.Video{}).Where("id = ? and uid = ?", video.ID, uid).Updates(
		map[string]interface{}{
			"title":        video.Title,
			"cover":        video.Cover,
			"introduction": video.Introduction,
			"original":     video.Original,
		},
	).Error; err != nil {
		tx.Rollback()
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = "视频信息修改失败"
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
		res.Msg = "更新审核状态失败"
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
		Msg:        "ok",
	}

	DB := common.GetDB()
	if !IsUserOwnsVideo(DB, vid, uid.(uint)) {
		//该视频不属于这个用户
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = "删除失败"
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
		Msg:        "ok",
	}
}

/*********************************************************
** 函数功能: 视频信息修改请求
** 日    期: 2021/11/10
**********************************************************/
func UpdateRequestService(review dto.UpdateVideoReviewRequest, uid interface{}) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}
	//从上下文中获取用户id
	DB := common.GetDB()
	tx := DB.Begin()
	if err := tx.Model(&model.Video{}).Where("id = ? and uid = ?", review.ID, uid).Updates(
		map[string]interface{}{
			"review": false,
		},
	).Error; err != nil {
		tx.Rollback()
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = "修改失败"
		return res
	}
	//更新审核状态
	if err := tx.Model(&model.Review{}).Where("vid = ?", review.ID).Updates(
		map[string]interface{}{
			"status": review.Status,
		},
	).Error; err != nil {
		tx.Rollback()
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = "状态更新失败"
		return res
	}
	tx.Commit()
	return res
}

/*********************************************************
** 函数功能: 通过视频ID获取视频
** 日    期: 2021/11/10
** 修改时间: 2021年11月17日12:59:45
** 版    本: 3.5.0
** 修改内容: 移除子视频
**********************************************************/
func GetVideoByIDService(vid int) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}

	var video model.Video
	DB := common.GetDB()
	DB.Model(&model.Video{}).Preload("Author").Where("id = ? and review = true", vid).First(&video)
	if video.ID == 0 {
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.CheckFailCode
		res.Msg = "视频不见了"
		return res
	}

	//视频数据
	like, collect := CollectAndLikeCount(DB, uint(vid))
	//增加播放量
	Redis := common.RedisClient
	if Redis != nil {
		strClicks, _ := Redis.Get(util.VideoClicksKey(vid)).Result()
		if strClicks == "" {
			Redis.RPush(util.ClicksVideoList, vid)
			//25小时防止数据当天过期
			Redis.Set(util.VideoClicksKey(vid), video.Clicks, time.Hour*25)
		}
		Redis.Incr(util.VideoClicksKey(vid))
	}
	var data = vo.VideoData{
		LikeCount:    like,
		CollectCount: collect,
	}
	res.Data = gin.H{"video": vo.ToVideoVo(video, data)}
	return res
}

/*********************************************************
** 函数功能: 获取视频交互数据
** 日    期:2021/11/11
**********************************************************/
func GetVideoInteractiveDataService(vid int, uid interface{}) response.ResponseStruct {
	DB := common.GetDB()
	//获取作者id
	var fid dto.AuthorUid
	DB.Raw("select uid from videos where id = ?", vid).Scan(&fid)
	like, collect := IsCollectAndLike(DB, uid.(uint), uint(vid))
	follow := IsFollow(DB, uid.(uint), fid.UID)
	data := vo.InteractiveData{
		Collect: collect,
		Like:    like,
		Follow:  follow,
	}

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"data": data},
		Msg:        "ok",
	}
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
		Msg:        "ok",
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
	Redis := common.RedisClient
	const sql = "select videos.id,title,cover,name as author,clicks from users,videos where users.id=videos.uid and review=1 and videos.deleted_at is null order by clicks desc"

	DB.Raw(sql).Scan(&videos)
	length := len(videos)
	//获取到播放量
	if Redis != nil {
		for i := 0; i < length; i++ {
			videos[i].Clicks = vo.GetClicksFromRedis(Redis, int(videos[i].ID), videos[i].Clicks)
		}
	}

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"videos": videos},
		Msg:        "ok",
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
		sql := "review = 1 and partition_id in (" + list + ")"
		Pagination.Model(&model.Video{}).Select("id,title,cover").Where(sql).Scan(&videos).Count(&total)
	}

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"count": total, "videos": videos},
		Msg:        "ok",
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
		Msg:        "ok",
	}
	var videos []vo.SearchVideoVo

	DB := common.GetDB()
	if !IsUserExist(DB, uint(uid)) {
		res.HttpStatus = http.StatusUnprocessableEntity
		res.Code = response.CheckFailCode
		res.Msg = "用户不存在"
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
** 函数功能: 通过视频ID获取子视频列表
** 日    期: 2021/11/11
** 修改时间: 2021年11月17日13:00:52
** 版    本: 3.5.0
** 修改内容: 移除子视频
**********************************************************/
func GetSubVideoListByVideoIDService(uid interface{}, page int, pageSize int, parentId int) response.ResponseStruct {
	//记录总数
	var totalSize int
	//分页查询
	var videos []model.Video
	DB := common.GetDB()
	DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
	DB.Where("uid = ?", uid, parentId).Find(&videos).Count(&totalSize)

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"count": totalSize, "data": vo.ToUploadVideoVo(videos)},
		Msg:        "ok",
	}
}

/*********************************************************
** 函数功能: 管理员获取视频列表
** 日    期: 2021年11月12日15:30:26
**********************************************************/
func AdminGetVideoListService(page int, pageSize int) response.ResponseStruct {
	var total int //记录总数
	var videos []model.Video

	DB := common.GetDB()
	DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
	DB.Where("review = 1").Find(&videos).Count(&total)

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"count": total, "videos": vo.ToAdminVideoVo(videos)},
		Msg:        "ok",
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
		Msg:        "ok",
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
		Msg:        "ok",
	}

	newVideo := model.Video{
		Title:        video.Title,
		Cover:        video.Cover,
		Introduction: video.Introduction,
		Original:     true,
		Uid:          0,
		VideoType:    "mp4",
		Video:        video.Video,
		Review:       true,
	}
	DB := common.GetDB()
	if err := DB.Create(&newVideo).Error; err != nil {
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.CheckFailCode
		res.Msg = "上传失败"
		return res
	}
	res.Data = gin.H{"vid": newVideo.ID}
	return res
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
	db.Where("id = ?", vid).First(&video)
	if video.ID != 0 {
		return true
	}
	return false
}

/*********************************************************
** 函数功能: 将点击量存入数据库
** 日    期:2021/7/22
**********************************************************/
func ClicksStoreInDB() {
	util.Logfile("[Info]", " Clicks are stored in the database")
	var vid int          //视频id
	var key string       //redis的key
	var clicks int       //点击量数字
	var strClicks string //字符串格式
	DB := common.GetDB()
	Redis := common.RedisClient
	if Redis == nil {
		util.Logfile("[Error]", " Clicks save failed")
		return
	}
	videos := Redis.LRange(util.ClicksVideoList, 0, -1).Val()
	for _, i := range videos {
		vid, _ = strconv.Atoi(i)
		key = util.VideoClicksKey(vid)
		strClicks, _ = Redis.Get(key).Result()
		clicks, _ = strconv.Atoi(strClicks)
		//删除redis数据
		Redis.Del(key)
		//写入数据库
		DB.Model(&model.Video{}).Where("id = ?", vid).Update("clicks", clicks)
	}
	//删除list
	Redis.Del(util.ClicksVideoList)
	util.Logfile("[Info]", " Click volume storage completed")
}
