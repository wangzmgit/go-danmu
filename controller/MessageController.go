package controller

import (
	"strconv"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"

	"github.com/gin-gonic/gin"
)

//发送私信
func SendMessage(ctx *gin.Context) {
	var requestMsg = model.Message{}
	err := ctx.Bind(&requestMsg)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	fid := requestMsg.Fid
	content := requestMsg.Content
	uid, _ := ctx.Get("id")
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
	DB := common.GetDB()
	DB.Create(&model.Message{Uid: uid.(uint), Fid: fid, FromId: uid.(uint), ToId: fid, Content: content})
	//切换消息归属人
	DB.Create(&model.Message{Uid: fid, Fid: uid.(uint), FromId: uid.(uint), ToId: fid, Content: content})
	response.Success(ctx, nil, "ok")
}

func GetMessageList(ctx *gin.Context) {
	DB := common.GetDB()
	//从上下文中获取用户id
	uid, _ := ctx.Get("id")
	var messageList []dto.MessagesListDto
	var sql = "select messages.id,messages.created_at,users.id as uid,users.name,users.avatar from messages,users "
	sql += "where messages.id in (select Max(id) from messages where deleted_at is null group by fid) and messages.fid = users.id and uid = ?"
	DB.Raw(sql, uid).Scan(&messageList)
	response.Success(ctx, gin.H{"messages": messageList}, "ok")
}

func GetMessageDetails(ctx *gin.Context) {
	DB := common.GetDB()
	uid, _ := ctx.Get("id")
	var messageDetails []dto.MessageDetailsDto
	fid, _ := strconv.Atoi(ctx.Query("fid"))
	if fid == 0 {
		response.Fail(ctx, nil, "消息不存在")
		return
	}
	DB.Model(&model.Message{}).Select("fid,from_id,content,created_at").Where("uid = ? AND fid = ?", uid.(uint), fid).Scan(&messageDetails)
	//查询用户信息
	var userInfo model.User
	DB.First(&userInfo, fid)
	response.Success(ctx, gin.H{"avatar": userInfo.Avatar, "name": userInfo.Name, "messages": messageDetails}, "ok")
}
