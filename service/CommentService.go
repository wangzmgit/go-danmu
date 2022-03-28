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
** 函数功能: 获取评论列表v1
** 日    期: 2021年11月11日17:51:55
**********************************************************/
func GetCommentsService(page int, pageSize int, vid int) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}
	var count int
	var comments []vo.CommentVo
	//查询语句
	sqlComment := "select comments.id,comments.created_at,content,uid,users.name,users.avatar from comments,users where comments.deleted_at is null and comments.uid = users.id and vid = ? limit ? offset ?"
	sqlReply := "select replies.id,replies.created_at,content,uid,users.name,users.avatar,reply_uid,reply_name from replies,users where replies.deleted_at is null and replies.uid = users.id and cid = ?"
	DB := common.GetDB()
	if !isVideoExist(DB, uint(vid)) {
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = response.VideoNotExist
		return res
	}

	DB.Model(&model.Comment{}).Where("vid = ?", vid).Count(&count)
	DB.Raw(sqlComment, vid, pageSize, (page-1)*pageSize).Scan(&comments)
	for i := 0; i < len(comments); i++ {
		//查询回复
		DB.Raw(sqlReply, comments[i].ID).Scan(&comments[i].Reply)
	}

	res.Data = gin.H{"count": count, "comments": comments}
	return res
}

/*********************************************************
** 函数功能: 获取评论列表v2
** 日    期: 2021年11月11日20:58:41
**********************************************************/
func GetCommentsV2Service(page int, pageSize int, vid int) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}
	var count int
	var comments []vo.CommentVo
	sqlComment := "select comments.id,comments.created_at,content,uid,users.name,users.avatar,reply_count " +
		"from comments,users where comments.deleted_at is null and comments.uid = users.id and vid = ?"
	sqlReply := "select content,users.name,reply_uid,reply_name from replies,users " +
		"where replies.deleted_at is null and replies.uid = users.id and cid = ? limit 2"

	DB := common.GetDB()

	if !isVideoExist(DB, uint(vid)) {
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = response.VideoNotExist
		return res
	}
	DB.Model(&model.Comment{}).Where("vid = ?", vid).Count(&count)
	DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
	DB.Raw(sqlComment, vid).Scan(&comments)
	for i := 0; i < len(comments); i++ {
		//查询回复
		DB.Raw(sqlReply, comments[i].ID).Scan(&comments[i].Reply)
	}
	res.Data = gin.H{"count": count, "comments": comments}
	return res
}

/*********************************************************
** 函数功能: 获取回复详情v2
** 日    期: 2021年11月11日21:04:11
**********************************************************/
func GetReplyDetailsV2Service(cid int, page int, pageSize int) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}

	var replies []vo.ReplyVo
	sql := "select replies.id,replies.created_at,content,uid,users.name,users.avatar,reply_uid,reply_name " +
		"from replies,users where replies.deleted_at is null and replies.uid = users.id and cid = ?"
	DB := common.GetDB()
	if !isCommentExist(DB, uint(cid)) {
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = response.CommentNotExist
		return res
	}
	DB = DB.Limit(pageSize).Offset((page - 1) * pageSize)
	DB.Model(&model.Reply{}).Where("cid = ?", cid)
	DB.Raw(sql, cid).Scan(&replies)
	res.Data = gin.H{"replies": replies}
	return res
}

/*********************************************************
** 函数功能: 删除评论
** 日    期: 2021年11月11日21:10:40
**********************************************************/
func DeleteCommentService(id uint, uid interface{}) response.ResponseStruct {
	DB := common.GetDB()
	DB.Where("id = ? and uid = ?", id, uid).Delete(model.Comment{})
	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}
}

/*********************************************************
** 函数功能: 删除回复
** 日    期: 2021年11月11日21:13:05
**********************************************************/
func DeleteReplyService(id uint, uid interface{}) response.ResponseStruct {
	DB := common.GetDB()
	DB.Where("id = ? and uid = ?", id, uid).Delete(model.Reply{})
	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}
}

/*********************************************************
** 函数功能: 评论
** 日    期:2021/7/27
**********************************************************/
func CommentService(comment dto.CommentDto, uid interface{}) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}

	DB := common.GetDB()
	if !isVideoExist(DB, comment.Vid) {
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = response.VideoNotExist
		return res
	}

	DB.Create(&model.Comment{Vid: comment.Vid, Content: comment.Content, Uid: uid.(uint)})
	return res
}

/*********************************************************
** 函数功能: 回复
** 日    期: 2021年11月11日21:20:59
**********************************************************/
func ReplyService(reply dto.ReplyDto, uid interface{}) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}

	DB := common.GetDB()
	if !isCommentExist(DB, reply.Cid) {
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = response.CommentNotExist
		return res
	}

	newReply := model.Reply{
		Cid:      reply.Cid,
		Content:  reply.Content,
		Uid:      uid.(uint),
		ReplyUid: reply.ReplyUid,
	}
	DB.Create(&newReply)
	//回复数+1,用于评论v2接口
	DB.Model(&model.Comment{}).Where("id = ?", reply.Cid).UpdateColumn("reply_count", gorm.Expr("reply_count + ?", 1))
	return res
}

/*********************************************************
** 函数功能: 评论是否存在
** 日    期:2021/7/27
**********************************************************/
func isCommentExist(db *gorm.DB, cid uint) bool {
	var comment model.Comment
	db.Where("id = ?", cid).First(&comment)
	if comment.ID != 0 {
		return true
	}
	return false
}
