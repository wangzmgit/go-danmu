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
	sql := "select messages.id,messages.created_at,`status`,users.id as uid,users.name,users.avatar from messages,users "
	sql += "where messages.id in (select Max(id) from messages where deleted_at is null group by fid)"
	sql += " and messages.fid = users.id and uid = ? order by id desc"
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
	DB.Model(&model.Message{}).Select("fid,from_id,content,created_at").
		Where("uid = ? AND fid = ?", uid.(uint), fid).Scan(&messageDetails)
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

/*********************************************************
** 函数功能: 获取消息详细信息V2
** 日    期: 2022年2月26日17:47:01
**********************************************************/
func GetMessageDetailsServiceV2(uid interface{}, fid, page, pageSize int) response.ResponseStruct {
	var userInfo model.User
	var messages []vo.MessageDetailsVo

	DB := common.GetDB()
	DB = DB.Limit(pageSize).Offset((page - 1) * pageSize).Order("id desc")
	DB.Model(&model.Message{}).Select("fid,from_id,content,created_at").Where("uid = ? AND fid = ?", uid.(uint), fid).Scan(&messages)
	// 此时查询到的消息为为倒叙，需要进行反转
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	//仅在获取第一页时查询用户信息
	if page == 1 {
		DB.First(&userInfo, fid)
	}

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"avatar": userInfo.Avatar, "name": userInfo.Name, "messages": messages},
		Msg:        response.OK,
	}
}

/*********************************************************
** 函数功能: 已读消息
** 日    期: 2022年3月3日19:40:20
**********************************************************/
func ReadMessageService(uid interface{}, fid uint) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        response.OK,
	}
	DB := common.GetDB()
	if err := DB.Model(&model.Message{}).
		Where("uid = ? and fid = ?", uid, fid).Update("status", 1).Error; err != nil {
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = response.UpdateStatusFail
		return res
	}

	return res
}
