package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/common"
	"kuukaa.fun/danmu-v4/model"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/util"
)

/*********************************************************
** 函数功能: 获取消息Websocket连接
** 日    期: 2022年2月25日09:31:10
**********************************************************/
func GetMsgConnect(ctx *gin.Context) {
	token := ctx.Query("token")
	auth, uid := AnalysisToken(token)
	if !auth {
		response.Fail(ctx, nil, response.PleaseLoginFirst)
		return
	}
	util.Logfile(util.InfoLog, " User "+uid+" connection successful")
	// 升级为websocket长链接
	common.WsHandler(ctx.Writer, ctx.Request, uid)
}

/*********************************************************
** 函数功能: 手动解析token
** 日    期: 2022年2月25日09:47:02
**********************************************************/
func AnalysisToken(tokenString string) (bool, string) {
	if tokenString == "" {
		return false, ""
	}

	token, claims, err := common.ParseUserToken(tokenString)
	if err != nil || !token.Valid {
		return false, ""
	}

	//通过token获取userID
	userId := claims.UserId
	DB := common.GetDB()
	var user model.User
	DB.First(&user, userId)

	if user.ID == 0 {
		return false, ""
	}
	return true, strconv.Itoa(int(userId))
}
