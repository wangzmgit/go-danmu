package service

import (
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"kuukaa.fun/danmu-v4/common"
	"kuukaa.fun/danmu-v4/model"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/util"
)

/*********************************************************
** 函数功能: 添加收藏
** 日    期:2021/11/11
**********************************************************/
func CollectService(vid uint, uid interface{}) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}
	//验证视频是否存在
	DB := common.GetDB()
	if !isVideoExist(DB, vid) {
		res.HttpStatus = http.StatusUnprocessableEntity
		res.Code = response.CheckFailCode
		res.Msg = response.VideoNotExist
		return res
	}
	//验证是否已经收藏
	status := isCollect(DB, uid.(uint), vid)
	if status == 0 {
		res.HttpStatus = http.StatusUnprocessableEntity
		res.Code = response.CheckFailCode
		res.Msg = response.IsCollect
		return res
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
	Redis := common.RedisClient
	if Redis != nil {
		strCollect, _ := Redis.Get(util.VideoCollectKey(intVid)).Result()
		if strCollect != "" {
			Redis.Incr(util.VideoCollectKey(intVid))
		}
	}
	return res
}

/*********************************************************
** 函数功能: 取消收藏
** 日    期:2021/11/11
**********************************************************/
func CancelCollectService(vid uint, uid interface{}) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}
	//验证收藏是否存在
	DB := common.GetDB()
	status := isCollect(DB, uid.(uint), vid)
	if status != 0 {
		res.HttpStatus = http.StatusUnprocessableEntity
		res.Code = response.CheckFailCode
		res.Msg = response.NotCollect
		return res
	} else {
		DB.Model(&model.Interactive{}).Where("uid = ? AND vid = ?", uid, vid).Update("collect", false)
		intVid := int(vid)
		Redis := common.RedisClient
		if Redis != nil {
			strCollect, _ := Redis.Get(util.VideoCollectKey(intVid)).Result()
			if strCollect != "" {
				Redis.Decr(util.VideoCollectKey(intVid))
			}
		}
		return res
	}
}

/*********************************************************
** 函数功能: 点赞
** 日    期:2021/11/11
**********************************************************/
func LikeService(vid uint, uid interface{}) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}
	//验证视频是否存在
	DB := common.GetDB()
	if !isVideoExist(DB, vid) {
		res.HttpStatus = http.StatusUnprocessableEntity
		res.Code = response.CheckFailCode
		res.Msg = response.VideoNotExist
		return res
	}
	//验证是否已经点赞
	status := isLike(DB, uid.(uint), vid)
	if status == 0 {
		res.HttpStatus = http.StatusUnprocessableEntity
		res.Code = response.CheckFailCode
		res.Msg = response.IsLike
		return res
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
	Redis := common.RedisClient
	if Redis != nil {
		strLike, _ := Redis.Get(util.VideoLikeKey(intVid)).Result()
		if strLike != "" {
			Redis.Incr(util.VideoLikeKey(intVid))
		}
	}
	return res
}

/*********************************************************
** 函数功能: 取消赞
** 日    期:2021/11/11
**********************************************************/
func DislikeService(vid uint, uid interface{}) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}
	//验证点赞是否存在
	DB := common.GetDB()
	status := isLike(DB, uid.(uint), vid)
	if status == 0 {
		DB.Model(&model.Interactive{}).Where("uid = ? AND vid = ?", uid, vid).Update("like", false)
		intVid := int(vid)
		Redis := common.RedisClient
		if Redis != nil {
			strLike, _ := Redis.Get(util.VideoLikeKey(intVid)).Result()
			if strLike != "" {
				Redis.Decr(util.VideoLikeKey(intVid))
			}
		}
		return res
	} else {
		res.HttpStatus = http.StatusUnprocessableEntity
		res.Code = response.CheckFailCode
		res.Msg = response.NotLike
		return res
	}
}

/*********************************************************
** 函数功能: 是否已经收藏
** 日    期: 2021/7/22
**********************************************************/
func isCollect(db *gorm.DB, uid uint, vid uint) int {
	//不存在返回-1，存在但是没有收藏返回1，已经收藏返回0
	var favorites = model.Interactive{}
	db.Where("uid = ? AND vid = ?", uid, vid).First(&favorites)
	if favorites.ID == 0 {
		return -1
	} else if !favorites.Collect {
		return 1
	}
	return 0
}

/*********************************************************
** 函数功能: 是否已经点赞
** 日    期: 2021/7/22
**********************************************************/
func isLike(db *gorm.DB, uid uint, vid uint) int {
	//不存在返回-1，存在但是没有点赞返回1，已经点赞返回0
	var like = model.Interactive{}
	db.Where("uid = ? AND vid = ?", uid, vid).First(&like)
	if like.ID == 0 {
		return -1
	} else if !like.Like {
		return 1
	}
	return 0
}

/*********************************************************
** 函数功能: 是否已经点赞和收藏
** 日    期:2021/7/22
** 返 回 值:是否点赞，是否收藏
**********************************************************/
func isCollectAndLike(db *gorm.DB, uid uint, vid uint) (bool, bool) {
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
** 修改时间: 2021/11/6
** 版    本: 3.3.0
** 修改内容: 如果redis出现问题,返回的点赞和收藏数都为0
**********************************************************/
func collectAndLikeCount(db *gorm.DB, vid uint) (int, int) {
	var like int
	var collect int
	intVid := int(vid)
	Redis := common.RedisClient
	if Redis == nil {
		return 0, 0
	}
	strLike, _ := Redis.Get(util.VideoLikeKey(intVid)).Result()
	strCollect, _ := Redis.Get(util.VideoCollectKey(intVid)).Result()
	if strLike == "" || strCollect == "" {
		//like和SQL的关键词冲突了，需要写成`like`
		db.Model(&model.Interactive{}).Where("vid = ? and `like` = 1", vid).Count(&like)
		db.Model(&model.Interactive{}).Where("vid = ? and collect = 1", vid).Count(&collect)
		//写入redis，设置6小时过期
		Redis.Set(util.VideoLikeKey(intVid), like, time.Hour*6)
		Redis.Set(util.VideoCollectKey(intVid), collect, time.Hour*6)
		return like, collect
	}
	like, _ = strconv.Atoi(strLike)
	collect, _ = strconv.Atoi(strCollect)
	return like, collect
}
