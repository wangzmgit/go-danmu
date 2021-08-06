package controller

import (
	"strconv"
	"time"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/util"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type IDRequest struct {
	ID uint
}

/*********************************************************
** 函数功能: 添加收藏
** 日    期:2021/7/22
**********************************************************/
func Collect(ctx *gin.Context) {
	//获取参数
	var request = IDRequest{}
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	vid := request.ID
	uid, _ := ctx.Get("id")
	//验证数据
	if vid <= 0 {
		response.CheckFail(ctx, nil, "视频不见了")
		return
	}
	//验证视频是否存在
	DB := common.GetDB()
	if !IsVideoExist(DB, vid) {
		response.CheckFail(ctx, nil, "视频不见了")
		return
	}
	//验证是否已经收藏
	status := IsCollect(DB, uid.(uint), vid)
	if status == 0 {
		response.CheckFail(ctx, nil, "已经收藏")
		return
	}
	if status == -1 {
		newFavorites := model.Interactive{
			Uid:     uid.(uint),
			Vid:     vid,
			Collect: true,
		}
		DB.Create(&newFavorites)
	} else {
		DB.Model(&model.Interactive{}).Where("uid = ? AND vid = ?", uid, vid).Update("collect", true)
	}
	intVid := int(vid)
	strCollect, _ := common.RedisClient.Get(util.VideoCollectKey(intVid)).Result()
	if strCollect != "" {
		common.RedisClient.Incr(util.VideoCollectKey(intVid))
	}
	response.Success(ctx, nil, "ok")
}

/*********************************************************
** 函数功能: 取消收藏
** 日    期:2021/7/22
**********************************************************/
func CancelCollect(ctx *gin.Context) {
	//获取参数
	var request = IDRequest{}
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	vid := request.ID
	uid, _ := ctx.Get("id")
	//验证收藏是否存在
	DB := common.GetDB()
	status := IsCollect(DB, uid.(uint), vid)
	if status != 0 {
		response.CheckFail(ctx, nil, "没有收藏")
	} else {
		DB.Model(&model.Interactive{}).Where("uid = ? AND vid = ?", uid, vid).Update("collect", false)
		intVid := int(vid)
		strCollect, _ := common.RedisClient.Get(util.VideoCollectKey(intVid)).Result()
		if strCollect != "" {
			common.RedisClient.Decr(util.VideoCollectKey(intVid))
		}
		response.Success(ctx, nil, "ok")
	}
}

/*********************************************************
** 函数功能: 点赞
** 日    期:2021/7/22
**********************************************************/
func Like(ctx *gin.Context) {
	//获取参数
	var request = IDRequest{}
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	vid := request.ID
	uid, _ := ctx.Get("id")
	//验证数据
	if vid <= 0 {
		response.CheckFail(ctx, nil, "视频不见了")
		return
	}
	//验证视频是否存在
	DB := common.GetDB()
	if !IsVideoExist(DB, vid) {
		response.CheckFail(ctx, nil, "视频不见了")
		return
	}
	//验证是否已经点赞
	status := IsLike(DB, uid.(uint), vid)
	if status == 0 {
		response.CheckFail(ctx, nil, "已经点过赞了")
		return
	}
	if status == -1 {
		newFavorites := model.Interactive{
			Uid:  uid.(uint),
			Vid:  vid,
			Like: true,
		}
		DB.Create(&newFavorites)
	} else {
		DB.Model(&model.Interactive{}).Where("uid = ? AND vid = ?", uid, vid).Update("like", true)
	}
	intVid := int(vid)
	strLike, _ := common.RedisClient.Get(util.VideoLikeKey(intVid)).Result()
	if strLike != "" {
		common.RedisClient.Incr(util.VideoLikeKey(intVid))
	}
	response.Success(ctx, nil, "ok")
}

/*********************************************************
** 函数功能: 取消赞
** 日    期:2021/7/22
**********************************************************/
func Dislike(ctx *gin.Context) {
	//获取参数
	var request = IDRequest{}
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	vid := request.ID
	uid, _ := ctx.Get("id")
	//验证点赞是否存在
	DB := common.GetDB()
	status := IsLike(DB, uid.(uint), vid)
	if status == 0 {
		DB.Model(&model.Interactive{}).Where("uid = ? AND vid = ?", uid, vid).Update("like", false)
		intVid := int(vid)
		strLike, _ := common.RedisClient.Get(util.VideoLikeKey(intVid)).Result()
		if strLike != "" {
			common.RedisClient.Decr(util.VideoLikeKey(intVid))
		}
		response.Success(ctx, nil, "ok")
	} else {
		response.CheckFail(ctx, nil, "还没有点赞")
	}
}

/*********************************************************
** 函数功能: 是否已经收藏
** 日    期:2021/7/22
**********************************************************/
func IsCollect(db *gorm.DB, uid uint, vid uint) int {
	//不存在返回-1，存在但是没有收藏返回1，已经收藏返回0
	var favorites = model.Interactive{}
	db.Where("uid = ? AND vid = ?", uid, vid).First(&favorites)
	if favorites.ID == 0 {
		return -1
	} else if favorites.Collect == false {
		return 1
	}
	return 0
}

/*********************************************************
** 函数功能: 是否已经点赞
** 日    期:2021/7/22
**********************************************************/
func IsLike(db *gorm.DB, uid uint, vid uint) int {
	//不存在返回-1，存在但是没有点赞返回1，已经点赞返回0
	var like = model.Interactive{}
	db.Where("uid = ? AND vid = ?", uid, vid).First(&like)
	if like.ID == 0 {
		return -1
	} else if like.Like == false {
		return 1
	}
	return 0
}

/*********************************************************
** 函数功能: 是否已经点赞和收藏
** 日    期:2021/7/22
** 返 回 值:是否点赞，是否收藏
**********************************************************/
func IsCollectAndLike(db *gorm.DB, uid uint, vid uint) (bool, bool) {
	var data = model.Interactive{}
	db.Where("uid = ? AND vid = ?", uid, vid).First(&data)
	if data.ID != 0 {
		return data.Like, data.Collect
	} else {
		return false, false
	}
}

/*********************************************************
** 函数功能: 点赞和收藏数据
** 日    期:2021/7/22
** 返 回 值:点赞数，收藏数
**********************************************************/
func CollectAndLikeCount(db *gorm.DB, vid uint) (int, int) {
	var like int
	var collect int
	intVid := int(vid)
	strLike, _ := common.RedisClient.Get(util.VideoLikeKey(intVid)).Result()
	strCollect, _ := common.RedisClient.Get(util.VideoCollectKey(intVid)).Result()
	if strLike == "" || strCollect == "" {
		//like和SQL的关键词冲突了，需要写成`like`
		db.Model(&model.Interactive{}).Where("vid = ? and `like` = 1", vid).Count(&like)
		db.Model(&model.Interactive{}).Where("vid = ? and collect = 1", vid).Count(&collect)
		//写入redis，设置6小时过期
		common.RedisClient.Set(util.VideoLikeKey(intVid), like, time.Hour*6)
		common.RedisClient.Set(util.VideoCollectKey(intVid), collect, time.Hour*6)
		return like, collect
	}
	like, _ = strconv.Atoi(strLike)
	collect, _ = strconv.Atoi(strCollect)
	return like, collect
}
