package controller

import (
	"net/http"
	"strconv"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/service"

	"github.com/gin-gonic/gin"
)

/*********************************************************
** 函数功能: 关注
** 日    期:2021/7/24
**********************************************************/
func Following(ctx *gin.Context) {
	//获取参数
	var follow dto.FollowRequest
	err := ctx.Bind(&follow)
	if err != nil {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "请求错误")
		return
	}
	//关注的人的id和自己的id
	fid := follow.ID
	uid, _ := ctx.Get("id")
	//判断关注的是否为自己
	if fid == uid {
		response.CheckFail(ctx, nil, "不能关注自己")
		return
	}

	res := service.FollowingService(fid, uid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 取消关注
** 日    期:2021/7/24
**********************************************************/
func UnFollow(ctx *gin.Context) {
	var follow dto.FollowRequest
	err := ctx.Bind(&follow)
	if err != nil {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "请求错误")
		return
	}
	//关注的人的id和自己的id
	fid := follow.ID
	uid, _ := ctx.Get("id")

	res := service.UnFollowService(fid, uid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 获取关注状态
** 日    期:2021/7/25
**********************************************************/
func GetFollowStatus(ctx *gin.Context) {
	fid, _ := strconv.Atoi(ctx.Query("fid"))
	if fid == 0 {
		response.CheckFail(ctx, nil, "用户不存在")
		return
	}
	uid, _ := ctx.Get("id")
	res := service.GetFollowStatusService(uid.(uint), uint(fid))
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 通过UID获取关注列表
** 日    期:2021/7/25
**********************************************************/
func GetFollowingByID(ctx *gin.Context) {
	uid, _ := strconv.Atoi(ctx.Query("uid"))
	if uid == 0 {
		response.CheckFail(ctx, nil, "用户不存在")
		return
	}

	res := service.GetFollowingByIDService(uid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 通过UID获取粉丝列表
** 日    期:2021/7/25
**********************************************************/
func GetFollowersByID(ctx *gin.Context) {
	uid, _ := strconv.Atoi(ctx.Query("uid"))
	if uid == 0 {
		response.CheckFail(ctx, nil, "用户不存在")
		return
	}

	res := service.GetFollowersByIDService(uid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 通过UID获取粉丝数
** 日    期:2021/7/25
**********************************************************/
func GetFollowCount(ctx *gin.Context) {
	uid, _ := strconv.Atoi(ctx.Query("uid"))
	if uid == 0 {
		response.CheckFail(ctx, nil, "用户不存在")
		return
	}

	res := service.GetFollowCountService(uid)
	response.HandleResponse(ctx, res)
}
