package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/common"
	"kuukaa.fun/danmu-v4/model"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/vo"
)

/*********************************************************
** 函数功能: 发送私信
** 日    期: 2021年11月12日11:18:49
**********************************************************/
func SendMessageService(uid uint, fid uint, content string) response.ResponseStruct {
	DB := common.GetDB()
	DB.Create(&model.Message{Uid: uid, Fid: fid, FromId: uid, ToId: fid, Content: content})
	// //切换消息归属人
	DB.Create(&model.Message{Uid: fid, Fid: uid, FromId: uid, ToId: fid, Content: content})

	//推送消息给接收者
	data, _ := json.Marshal(&vo.MessageVo{
		Fid:     uid,
		Content: content,
	})
	common.SendMsgToUser(strconv.Itoa(int(fid)), data)

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}
}

/*********************************************************
** 函数功能: 获取消息列表
** 日    期: 2021年11月12日11:21:52
**********************************************************/
func GetMessageListService(uid interface{}) response.ResponseStruct {
	DB := common.GetDB()
	var messageList []vo.MessagesListVo
	var sql = "select messages.id,messages.created_at,users.id as uid,users.name,users.avatar from messages,users "
	sql += "where messages.id in (select Max(id) from messages where deleted_at is null group by fid) and messages.fid = users.id and uid = ?"
	DB.Raw(sql, uid).Scan(&messageList)

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"messages": messageList},
		Msg:        response.OK,
	}
}

/*********************************************************
** 函数功能: 获取消息详细信息
** 日    期: 2021年11月12日11:23:56
**********************************************************/
func GetMessageDetailsService(uid interface{}, fid int) response.ResponseStruct {
	var messageDetails []vo.MessageDetailsVo

	DB := common.GetDB()
	DB.Model(&model.Message{}).Select("fid,from_id,content,created_at").Where("uid = ? AND fid = ?", uid.(uint), fid).Scan(&messageDetails)
	//查询用户信息
	var userInfo model.User
	DB.First(&userInfo, fid)

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"avatar": userInfo.Avatar, "name": userInfo.Name, "messages": messageDetails},
		Msg:        response.OK,
	}
}
