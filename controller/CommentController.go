package controller

import (
	"strconv"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func GetComments(ctx *gin.Context) {
	var comments []dto.CommentDto
	var count int
	//获取分页信息
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	vid, _ := strconv.Atoi(ctx.Query("vid"))
	sqlComment := "select comments.id,comments.created_at,content,uid,users.name,users.avatar from comments,users where comments.deleted_at is null and comments.uid = users.id and vid = ? limit ? offset ?"
	sqlReply := "select replies.id,replies.created_at,content,uid,users.name,users.avatar,reply_uid,reply_name from replies,users where replies.deleted_at is null and replies.uid = users.id and cid = ?"
	DB := common.GetDB()
	if !IsVideoExist(DB, uint(vid)) {
		response.Fail(ctx, nil, "视频不存在")
		return
	}
	if page > 0 && pageSize > 0 {
		DB.Model(&model.Comment{}).Where("vid = ?", vid).Count(&count)
		DB.Raw(sqlComment, vid, pageSize, (page-1)*pageSize).Scan(&comments)
		for i := 0; i < len(comments); i++ {
			//查询回复
			DB.Raw(sqlReply, comments[i].ID).Scan(&comments[i].Reply)
		}
		response.Success(ctx, gin.H{"count": count, "comments": comments}, "ok")
	} else {
		response.Fail(ctx, nil, "获取失败")
	}
}

/*********************************************************
** 函数功能: 删除评论
** 日    期:2021/7/27
**********************************************************/
func DeleteComment(ctx *gin.Context) {
	//获取参数
	var request = IDRequest{}
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	id := request.ID
	uid, _ := ctx.Get("id")
	DB := common.GetDB()
	DB.Where("id = ? and uid = ?", id, uid).Delete(model.Comment{})
	response.Success(ctx, nil, "ok")
}

/*********************************************************
** 函数功能: 删除回复
** 日    期:2021/7/27
**********************************************************/
func DeleteReply(ctx *gin.Context) {
	var request = IDRequest{}
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	id := request.ID
	uid, _ := ctx.Get("id")
	DB := common.GetDB()
	DB.Where("id = ? and uid = ?", id, uid).Delete(model.Reply{})
	response.Success(ctx, nil, "ok")
}

/*********************************************************
** 函数功能: 评论
** 日    期:2021/7/27
**********************************************************/
func Comment(ctx *gin.Context) {
	var comment = model.Comment{}
	err := ctx.Bind(&comment)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	content := comment.Content
	vid := comment.Vid
	uid, _ := ctx.Get("id")
	DB := common.GetDB()
	if !IsVideoExist(DB, vid) {
		response.CheckFail(ctx, nil, "视频不存在")
		return
	}
	if len(content) == 0 {
		response.CheckFail(ctx, nil, "评论不能为空")
		return
	}
	DB.Create(&model.Comment{Vid: vid, Content: content, Uid: uid.(uint)})
	response.Success(ctx, nil, "ok")
}

/*********************************************************
** 函数功能: 回复
** 日    期:2021/7/27
**********************************************************/
func Reply(ctx *gin.Context) {
	var reply = model.Reply{}
	err := ctx.Bind(&reply)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	cid := reply.Cid
	content := reply.Content
	replyUid := reply.ReplyUid
	replyName := reply.ReplyName
	uid, _ := ctx.Get("id")
	DB := common.GetDB()
	if cid == 0 || !IsCommentExist(DB, cid) {
		response.CheckFail(ctx, nil, "评论不存在")
		return
	}
	if len(content) == 0 {
		response.CheckFail(ctx, nil, "回复不能为空")
		return
	}
	newReply := model.Reply{
		Cid:       cid,
		Content:   content,
		Uid:       uid.(uint),
		ReplyUid:  replyUid,
		ReplyName: replyName,
	}
	DB.Create(&newReply)
	response.Success(ctx, nil, "ok")
}

/*********************************************************
** 函数功能: 评论是否存在
** 日    期:2021/7/27
**********************************************************/
func IsCommentExist(db *gorm.DB, cid uint) bool {
	var comment model.Comment
	db.Where("id = ?", cid).First(&comment)
	if comment.ID != 0 {
		return true
	}
	return false
}
