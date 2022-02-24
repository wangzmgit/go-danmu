package manage

import (
	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/dto"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/service"
)

/*********************************************************
** 函数功能: 删除合集
** 日    期: 2022年2月24日12:42:410
**********************************************************/
func AdminDeleteCollection(ctx *gin.Context) {
	var request dto.AdminIdDto
	if err := ctx.Bind(&request); err != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}
	id := request.ID

	res := service.AdminDeleteCollectionService(id)
	response.HandleResponse(ctx, res)
}
