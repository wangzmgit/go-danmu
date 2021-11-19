package service

import (
	"net/http"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"
)

func CreateCollectionService(collection dto.CreateCollectionDto, uid interface{}) response.ResponseStruct {
	DB := common.GetDB()

	newCollection := model.Collection{
		Title: collection.Title,
		Cover: collection.Cover,
		Desc:  collection.Desc,
		Uid:   uid.(uint),
	}

	DB.Create(&newCollection)
	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}
}
