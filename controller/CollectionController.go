package controller

import (
	"strconv"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/service"

	"github.com/gin-gonic/gin"
)

/*********************************************************
** 函数功能: 创建合集
** 日    期: 2021年11月19日14:48:29
** 版    本: 3.6.0
**********************************************************/
func CreateCollection(ctx *gin.Context) {
	var request dto.CreateCollectionDto
	err := ctx.Bind(&request)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	title := request.Title
	cover := request.Cover
	uid, _ := ctx.Get("id")

	//验证数据
	if len(title) == 0 {
		response.CheckFail(ctx, nil, "标题不能为空")
		return
	}

	if len(cover) == 0 {
		response.CheckFail(ctx, nil, "封面图不能为空")
		return
	}

	res := service.CreateCollectionService(request, uid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 修改合集信息
** 日    期: 2021年11月19日14:48:29
** 版    本: 3.6.0
**********************************************************/
func ModifyCollection(ctx *gin.Context) {
	var request dto.ModifyCollectionDto
	err := ctx.Bind(&request)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	title := request.Title
	cover := request.Cover
	uid, _ := ctx.Get("id")

	//验证数据
	if len(title) == 0 {
		response.CheckFail(ctx, nil, "标题不能为空")
		return
	}

	if len(cover) == 0 {
		response.CheckFail(ctx, nil, "封面图不能为空")
		return
	}

	res := service.ModifyCollectionService(request, uid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 获取自己创建的合集
** 日    期: 2021年11月19日19:58:43
** 版    本: 3.6.0
**********************************************************/
func GetCreateCollectionList(ctx *gin.Context) {
	uid, _ := ctx.Get("id")
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))

	if page <= 0 || pageSize <= 0 {
		response.CheckFail(ctx, nil, "页码或数量有误")
		return
	}

	if pageSize >= 30 {
		response.CheckFail(ctx, nil, "请求数量过多")
		return
	}

	res := service.GetCreateCollectionListService(page, pageSize, uid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 获取合集内容
** 日    期: 2021年11月19日20:01:08
** 版    本: 3.6.0
**********************************************************/
func GetCollectionContent(ctx *gin.Context) {
	cid, _ := strconv.Atoi(ctx.Query("cid"))
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))

	if cid <= 0 {
		response.CheckFail(ctx, nil, "视频不存在")
		return
	}

	if page <= 0 || pageSize <= 0 {
		response.CheckFail(ctx, nil, "页码或数量有误")
		return
	}

	if pageSize >= 30 {
		response.CheckFail(ctx, nil, "请求数量过多")
		return
	}

	res := service.GetCollectionContentService(cid, page, pageSize)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 获取合集信息
** 日    期: 2021年11月20日16:23:57
** 版    本: 3.6.0
**********************************************************/
func GetCollectionByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Query("id"))
	if id <= 0 {
		response.CheckFail(ctx, nil, "找不到合集")
		return
	}

	res := service.GetCollectionByIDService(id)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 删除合集
** 日    期: 2021年11月20日16:55:39
** 版    本: 3.6.0
**********************************************************/
func DeleteCollection(ctx *gin.Context) {
	//获取参数
	var collection dto.DeleteCollectionDto
	err := ctx.Bind(&collection)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	id := collection.ID
	uid, _ := ctx.Get("id")

	//数据验证
	if id == 0 {
		response.CheckFail(ctx, nil, "合集不存在")
		return
	}

	//删除合集
	res := service.DeleteCollectionService(id, uid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 获取可以添加的视频列表
** 日    期: 2021年11月20日20:38:13
** 版    本: 3.6.0
**********************************************************/
func GetCanAddVideo(ctx *gin.Context) {
	//获取参数
	uid, _ := ctx.Get("id")
	id, _ := strconv.Atoi(ctx.Query("id"))
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	if page <= 0 || pageSize <= 0 {
		response.CheckFail(ctx, nil, "页码或数量有误")
		return
	}

	if pageSize >= 30 {
		response.CheckFail(ctx, nil, "请求数量过多")
		return
	}

	if id <= 0 {
		response.CheckFail(ctx, nil, "找不到合集")
		return
	}

	//删除合集
	res := service.GetCanAddVideoService(id, uid, page, pageSize)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 向合集内添加视频
** 日    期: 2021年11月21日09:31:09
** 版    本: 3.6.0
**********************************************************/
func AddVideoToCollection(ctx *gin.Context) {
	//获取参数
	var video dto.AddVideoDto
	err := ctx.Bind(&video)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	vid := video.Vid
	cid := video.Cid
	uid, _ := ctx.Get("id")

	//数据验证
	if vid == 0 || cid == 0 {
		response.CheckFail(ctx, nil, "合集或视频不存在")
		return
	}

	//添加合集
	res := service.AddVideoToCollectionService(vid, cid, uid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 删除合集内的视频
** 日    期: 2021年11月21日13:32:59
** 版    本: 3.6.0
**********************************************************/
func DeleteCollectionVideo(ctx *gin.Context) {
	//获取参数
	var video dto.DeleteVideoDto
	err := ctx.Bind(&video)
	if err != nil {
		response.Fail(ctx, nil, "请求错误")
		return
	}
	vid := video.Vid
	cid := video.Cid
	uid, _ := ctx.Get("id")

	//数据验证
	if vid == 0 || cid == 0 {
		response.CheckFail(ctx, nil, "合集或视频不存在")
		return
	}

	//添加合集
	res := service.DeleteCollectionVideoService(vid, cid, uid)
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 获取合集列表
** 日    期: 2021年11月21日14:53:56
** 版    本: 3.6.0
**********************************************************/
func GetCollectionList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))

	if page <= 0 || pageSize <= 0 {
		response.CheckFail(ctx, nil, "页码或数量有误")
		return
	}

	if pageSize >= 30 {
		response.CheckFail(ctx, nil, "请求数量过多")
		return
	}

	res := service.GetCollectionListService(page, pageSize)
	response.HandleResponse(ctx, res)
}
