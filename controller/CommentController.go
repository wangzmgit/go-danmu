package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/dto"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/service"
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
		response.CheckFail(ctx, nil, response.PageOrSizeError)
		return
	}
	if vid <= 0 {
		response.CheckFail(ctx, nil, response.VideoNotExist)
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
		response.CheckFail(ctx, nil, response.PageOrSizeError)
		return
	}
	if pageSize >= 30 {
		response.CheckFail(ctx, nil, response.TooManyRequests)
		return
	}

	res := service.GetCommentsV2Service(page, pageSize, vid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 获取回复详情v2
** 日    期: 2021/10/5
**********************************************************/
func GetReplyDetailsV2(ctx *gin.Context) {
	//获取分页信息
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	cid, _ := strconv.Atoi(ctx.Query("cid"))
	if cid <= 0 {
		response.CheckFail(ctx, nil, response.ParameterError)
		return
	}
	if page <= 0 || pageSize <= 0 {
		response.CheckFail(ctx, nil, response.PageOrSizeError)
		return
	}
	if pageSize >= 30 {
		response.CheckFail(ctx, nil, response.TooManyRequests)
		return
	}

	res := service.GetReplyDetailsV2Service(cid, page, pageSize)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 删除评论
** 日    期:2021/7/27
**********************************************************/
func DeleteComment(ctx *gin.Context) {
	//获取参数
	var request dto.CommentIdDto
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, response.RequestError)
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
	var request dto.CommentIdDto
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, response.RequestError)
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
	var comment dto.CommentDto
	err := ctx.Bind(&comment)
	if err != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}
	content := comment.Content
	uid, _ := ctx.Get("id")

	if len(content) == 0 {
		response.CheckFail(ctx, nil, response.CommentCheck)
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
	var reply dto.ReplyDto
	err := ctx.Bind(&reply)
	if err != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}
	cid := reply.Cid
	content := reply.Content
	uid, _ := ctx.Get("id")

	if cid == 0 {
		response.CheckFail(ctx, nil, response.CommentNotExist)
		return
	}
	if len(content) == 0 {
		response.CheckFail(ctx, nil, response.CommentCheck)
		return
	}

	res := service.ReplyService(reply, uid)
	response.HandleResponse(ctx, res)
}
