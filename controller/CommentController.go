package controller

import (
	"strconv"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/service"

	"github.com/gin-gonic/gin"
)

/*********************************************************
** 函数功能: 获取评论列表
** 日    期:
**********************************************************/
func GetComments(ctx *gin.Context) {
	//获取分页信息
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	vid, _ := strconv.Atoi(ctx.Query("vid"))
	if page <= 0 || pageSize <= 0 {
		response.CheckFail(ctx, nil, "页码或数量有误")
		return
	}
	if vid <= 0 {
		response.CheckFail(ctx, nil, "视频不存在")
		return
	}

	res := service.GetCommentsService(page, pageSize, vid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 获取评论列表v2
** 日    期:2021/10/5
**********************************************************/
func GetCommentsV2(ctx *gin.Context) {

	//获取分页信息
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	vid, _ := strconv.Atoi(ctx.Query("vid"))

	if page <= 0 || pageSize <= 0 {
		response.CheckFail(ctx, nil, "页码或数量有误")
		return
	}
	res := service.GetCommentsV2Service(page, pageSize, vid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 获取回复详情v2
** 日    期:2021/10/5
**********************************************************/
func GetReplyDetailsV2(ctx *gin.Context) {
	//获取分页信息
	cid, _ := strconv.Atoi(ctx.Query("cid"))
	if cid <= 0 {
		response.CheckFail(ctx, nil, "参数有误")
		return
	}
	res := service.GetReplyDetailsV2Service(cid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 删除评论
** 日    期:2021/7/27
**********************************************************/
func DeleteComment(ctx *gin.Context) {
	//获取参数
	var request dto.CommentDeleteRequest
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	id := request.ID
	uid, _ := ctx.Get("id")

	res := service.DeleteCommentService(id, uid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 删除回复
** 日    期:2021/7/27
**********************************************************/
func DeleteReply(ctx *gin.Context) {
	var request dto.CommentDeleteRequest
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	id := request.ID
	uid, _ := ctx.Get("id")

	res := service.DeleteReplyService(id, uid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 评论
** 日    期:2021/7/27
**********************************************************/
func Comment(ctx *gin.Context) {
	var comment dto.CommentRequest
	err := ctx.Bind(&comment)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	content := comment.Content
	uid, _ := ctx.Get("id")

	if len(content) == 0 {
		response.CheckFail(ctx, nil, "评论不能为空")
		return
	}

	res := service.CommentService(comment, uid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 回复
** 日    期:2021/7/27
**********************************************************/
func Reply(ctx *gin.Context) {
	var reply dto.ReplyRequest
	err := ctx.Bind(&reply)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	cid := reply.Cid
	content := reply.Content
	uid, _ := ctx.Get("id")

	if cid == 0 {
		response.CheckFail(ctx, nil, "评论不存在")
		return
	}
	if len(content) == 0 {
		response.CheckFail(ctx, nil, "回复不能为空")
		return
	}

	res := service.ReplyService(reply, uid)
	response.HandleResponse(ctx, res)
}
