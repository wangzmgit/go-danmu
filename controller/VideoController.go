package controller

import (
	"strconv"
	"time"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/util"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
)

type Video struct {
	Title        string `json:"title"`
	Cover        string `json:"cover"`
	Introduction string `json:"introduction"`
}

//视频交互数据
type InteractiveData struct {
	Collect bool `json:"collect"`
	Like    bool `json:"like"`
	Follow  bool `json:"follow"`
}

//查询uid
type SelectUID struct {
	UID uint
}

/*********************************************************
** 函数功能: 上传视频信息
** 日    期:2021/7/16
**********************************************************/
func UploadVideoInfo(ctx *gin.Context) {
	var video = model.Video{}
	err := ctx.Bind(&video)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	title := video.Title
	cover := video.Cover
	introduction := video.Introduction
	original := video.Original
	//验证数据
	if len(title) == 0 {
		response.CheckFail(ctx, nil, "标题不能为空")
		return
	}
	if len(cover) == 0 {
		response.CheckFail(ctx, nil, "封面图不能为空")
		return
	}

	uid, _ := ctx.Get("id")
	newVideo := model.Video{
		Title:        title,
		Cover:        cover,
		Introduction: introduction,
		Uid:          uid.(uint),
		Original:     original,
	}
	DB := common.GetDB()
	tx := DB.Begin()
	if err := tx.Create(&newVideo).Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "上传失败")
		return
	}
	//创建审核状态
	if err := tx.Create(&model.Review{Vid: newVideo.ID, Status: 500}).Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "上传失败")
		return
	}
	tx.Commit()
	response.Success(ctx, gin.H{"vid": newVideo.ID}, "ok")
}

/*********************************************************
** 函数功能: 获取视频状态
** 日    期:2021/7/16
**********************************************************/
func GetVideoStatus(ctx *gin.Context) {
	var review = model.Review{}
	vid, _ := strconv.Atoi(ctx.Query("vid"))
	uid, _ := ctx.Get("id")
	DB := common.GetDB()
	DB.Model(&model.Review{}).Preload("Video").Where("vid = ?", vid).First(&review)
	if review.ID == 0 || review.Video.Uid != uid {
		response.Fail(ctx, nil, "视频不见了")
		return
	}
	var video = Video{
		Title:        review.Video.Title,
		Cover:        review.Video.Cover,
		Introduction: review.Video.Introduction,
	}
	response.Success(ctx, gin.H{"status": review.Status, "remarks": review.Remarks, "video": video}, "ok")
}

/*********************************************************
** 函数功能: 修改视频信息
** 日    期:2021/7/17
**********************************************************/
func ModifyVideoInfo(ctx *gin.Context) {
	//获取参数
	var video = model.Video{}
	err := ctx.Bind(&video)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	id := video.ID
	title := video.Title
	cover := video.Cover
	introduction := video.Introduction
	original := video.Original

	if len(title) == 0 {
		response.CheckFail(ctx, nil, "标题不能为空")
		return
	}
	if len(cover) == 0 {
		response.CheckFail(ctx, nil, "封面图不能为空")
		return
	}

	//从上下文中获取用户id
	uid, _ := ctx.Get("id")
	DB := common.GetDB()
	tx := DB.Begin()
	if err := tx.Model(&model.Video{}).Where("id = ? and uid = ?", id, uid).Updates(map[string]interface{}{"title": title, "cover": cover, "introduction": introduction, "original": original}).Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "修改失败")
		return
	}
	//更新审核状态
	if err := tx.Model(&model.Review{}).Where("vid = ?", id).Updates(map[string]interface{}{"status": 1000}).Error; err != nil {
		tx.Rollback()
		response.Fail(ctx, nil, "修改失败")
		return
	}
	tx.Commit()
	response.Success(ctx, nil, "ok")
}

/*********************************************************
** 函数功能: 删除视频
** 日    期:2021/7/17
**********************************************************/
func DeleteVideo(ctx *gin.Context) {
	//获取参数
	var video = model.Video{}
	err := ctx.Bind(&video)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	id := video.ID
	uid, _ := ctx.Get("id")
	DB := common.GetDB()
	DB.Where("id = ? and uid = ?", id, uid).Delete(model.Video{})
	//删除播放量数据
	common.RedisClient.Del(util.VideoClicksKey(int(id)))
	response.Success(ctx, nil, "ok")
}

/*********************************************************
** 函数功能: 获取自己的视频
** 日    期:2021/7/17
**********************************************************/
func GetMyUploadVideo(ctx *gin.Context) {
	//获取分页信息
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	if page > 0 && pageSize > 0 {
		//记录总数
		var totalSize int
		//分页查询
		var videos []model.Video
		uid, _ := ctx.Get("id")
		DB := common.GetDB()
		DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
		DB.Where("uid = ?", uid).Find(&videos).Count(&totalSize)
		response.Success(ctx, gin.H{"count": totalSize, "data": dto.ToUploadVideoDto(videos)}, "ok")
	} else {
		response.Fail(ctx, nil, "获取失败")
	}
}

/*********************************************************
** 函数功能: 视频信息修改请求
** 日    期:2021/7/18
**********************************************************/
func UpdateRequest(ctx *gin.Context) {
	var review = model.Review{}
	err := ctx.Bind(&review)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	id := review.ID
	status := review.Status
	if status == 5001 || status == 5002 {
		//从上下文中获取用户id
		uid, _ := ctx.Get("id")
		DB := common.GetDB()
		tx := DB.Begin()
		if err := tx.Model(&model.Video{}).Where("id = ? and uid = ?", id, uid).Updates(map[string]interface{}{"review": false}).Error; err != nil {
			tx.Rollback()
			response.Fail(ctx, nil, "修改失败")
			return
		}
		//更新审核状态
		if err := tx.Model(&model.Review{}).Where("vid = ?", id).Updates(map[string]interface{}{"status": status}).Error; err != nil {
			tx.Rollback()
			response.Fail(ctx, nil, "状态更新失败")
			return
		}
		tx.Commit()
		response.Success(ctx, nil, "ok")
	} else {
		response.Fail(ctx, nil, "申请状态有误")
	}
}

/*********************************************************
** 函数功能: 通过ID获取视频
** 日    期:2021/7/19
**********************************************************/
func GetVideoByID(ctx *gin.Context) {
	var video model.Video
	vid, _ := strconv.Atoi(ctx.Query("vid"))
	if vid == 0 {
		response.CheckFail(ctx, nil, "视频不见了")
		return
	}
	DB := common.GetDB()
	DB.Model(&model.Video{}).Preload("Author").Where("id = ? and review = true", vid).First(&video)
	if video.ID == 0 {
		response.CheckFail(ctx, nil, "视频不见了")
	} else {
		//视频数据
		like, collect := CollectAndLikeCount(DB, uint(vid))
		//增加播放量
		strClicks, _ := common.RedisClient.Get(util.VideoClicksKey(vid)).Result()
		if strClicks == "" {
			common.RedisClient.RPush(util.ClicksVideoList, vid)
			//25小时防止数据当天过期
			common.RedisClient.Set(util.VideoClicksKey(vid), video.Clicks, time.Hour*25)
		}
		common.RedisClient.Incr(util.VideoClicksKey(vid))
		var data = dto.VideoData{
			LikeCount:    like,
			CollectCount: collect,
		}
		response.Success(ctx, gin.H{"video": dto.ToVideoDto(video, data)}, "ok")
	}
}

/*********************************************************
** 函数功能: 获取视频交互数据
** 日    期:2021/7/22
**********************************************************/
func GetVideoInteractiveData(ctx *gin.Context) {
	vid, _ := strconv.Atoi(ctx.Query("vid"))
	if vid == 0 {
		response.CheckFail(ctx, nil, "视频不见了")
		return
	}
	DB := common.GetDB()
	uid, _ := ctx.Get("id")
	//获取作者id
	var fid SelectUID
	DB.Raw("select uid from videos where id = ?", vid).Scan(&fid)
	like, collect := IsCollectAndLike(DB, uid.(uint), uint(vid))
	follow := IsFollow(DB, uid.(uint), fid.UID)
	data := InteractiveData{
		Collect: collect,
		Like:    like,
		Follow:  follow,
	}
	response.Success(ctx, gin.H{"data": data}, "ok")
}

/*********************************************************
** 函数功能: 获取收藏列表
** 日    期:2021/7/22
**********************************************************/
func GetCollectVideo(ctx *gin.Context) {
	var favorites []model.Interactive
	var count int
	uid, _ := ctx.Get("id")
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	if page > 0 && pageSize > 0 {
		DB := common.GetDB()
		DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
		DB.Where("uid = ? AND collect = true", uid).Preload("Video").Find(&favorites).Count(&count)
		response.Success(ctx, gin.H{"count": count, "videos": dto.ToCollectVideoDto(favorites)}, "ok")
	} else {
		response.Fail(ctx, nil, "请求错误")
	}
}

/*********************************************************
** 函数功能: 获取推荐视频
** 日    期:2021/8/1
**********************************************************/
func GetRecommendVideo(ctx *gin.Context) {
	//因为视频比较少，就直接按播放量排名
	DB := common.GetDB()
	var videos []dto.RecommendVideo
	DB = DB.Limit(8)
	Redis := common.RedisClient
	const sql = "select videos.id,title,cover,name as author,clicks from users,videos where users.id=videos.uid and review=1 and videos.deleted_at is null order by clicks desc"
	DB.Raw(sql).Scan(&videos)
	length := len(videos)
	//获取到播放量
	for i := 0; i < length; i++ {
		videos[i].Clicks = dto.GetClicksFromRedis(Redis, int(videos[i].ID), videos[i].Clicks)
	}
	response.Success(ctx, gin.H{"videos": videos}, "ok")
}

/*********************************************************
** 函数功能: 获取视频列表
** 日    期:2021/8/1
**********************************************************/
func GetVideoList(ctx *gin.Context) {
	DB := common.GetDB()
	var videos []dto.SearchVideoDto
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	if page > 0 && pageSize > 0 {
		//记录总数
		var total int
		DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
		DB.Model(&model.Video{}).Select("id,title,cover").Where("review = 1").Scan(&videos).Count(&total)
		response.Success(ctx, gin.H{"count": total, "videos": videos}, "ok")
	} else {
		response.Fail(ctx, nil, "获取数量有误")
	}
}

/*********************************************************
** 函数功能: 通过用户ID获取视频列表
** 日    期:2021/8/4
**********************************************************/
func GetVideoListByUserID(ctx *gin.Context) {
	DB := common.GetDB()
	var videos []dto.SearchVideoDto
	uid, _ := strconv.Atoi(ctx.Query("uid"))
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	if !IsUserExist(DB, uint(uid)) {
		response.CheckFail(ctx, nil, "用户不存在")
		return
	}
	if page > 0 && pageSize > 0 {
		//记录总数
		var total int
		DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
		DB.Model(&model.Video{}).Select("id,title,cover").Where("review = 1 and uid = ?", uid).Scan(&videos).Count(&total)
		response.Success(ctx, gin.H{"count": total, "videos": videos}, "ok")
	} else {
		response.Fail(ctx, nil, "获取数量有误")
	}
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
	util.Logfile("[info]", " Clicks are stored in the database")
	var vid int          //视频id
	var key string       //redis的key
	var clicks int       //点击量数字
	var strClicks string //字符串格式
	DB := common.GetDB()
	Redis := common.RedisClient
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
	util.Logfile("[info]", " Click volume storage completed")
}
