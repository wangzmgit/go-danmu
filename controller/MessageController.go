package controller

import (
	"strconv"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/service"

	"github.com/gin-gonic/gin"
)

/*********************************************************
** 函数功能: 发送私信
**********************************************************/
func SendMessage(ctx *gin.Context) {
	var requestMsg dto.SendMessageReques
	err := ctx.Bind(&requestMsg)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	fid := requestMsg.Fid
	content := requestMsg.Content
	uid, _ := ctx.Get("id")
	//验证数据
	if fid == 0 {
		response.CheckFail(ctx, nil, "发送失败")
		return
	}
	if fid == uid.(uint) {
		response.CheckFail(ctx, nil, "不能发送给自己")
		return
	}
	if len(content) == 0 {
		response.CheckFail(ctx, nil, "不能发送空内容")
		return
	}

	res := service.SendMessageService(uid.(uint), fid, content)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 获取消息列表
**********************************************************/
func GetMessageList(ctx *gin.Context) {
	//从上下文中获取用户id
	uid, _ := ctx.Get("id")
	res := service.GetMessageListService(uid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 获取消息详细信息
**********************************************************/
func GetMessageDetails(ctx *gin.Context) {
	uid, _ := ctx.Get("id")
	fid, _ := strconv.Atoi(ctx.Query("fid"))
	if fid == 0 {
		response.Fail(ctx, nil, "消息不存在")
		return
	}

	res := service.GetMessageDetailsService(uid, fid)
	response.HandleResponse(ctx, res)
}
