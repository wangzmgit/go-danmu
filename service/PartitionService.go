package service

import (
	"net/http"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/dto"
	"wzm/danmu3.0/model"
	"wzm/danmu3.0/response"
	"wzm/danmu3.0/vo"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

/*********************************************************
** 函数功能: 获取分区
** 日    期: 2021年12月9日
**********************************************************/
func GetPartitionListService(fid int) response.ResponseStruct {
	var partitions []vo.PartitionVo

	DB := common.GetDB()
	DB.Model(&model.Partition{}).Select("id,content").Where("fid = ?", fid).Scan(&partitions)

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"partitions": partitions},
		Msg:        "ok",
	}
}

/*********************************************************
** 函数功能: 添加分区
** 日    期: 2021年12月9日
**********************************************************/
func AddPartitionService(partition dto.PartitionDto) response.ResponseStruct {
	res := response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}
	DB := common.GetDB()

	if partition.Fid != 0 && !IsParentPartitionExist(DB, partition.Fid) {
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = "所属分区不存在"
		return res
	}

	newPartition := model.Partition{
		Content: partition.Content,
		Fid:     partition.Fid,
	}

	if err := DB.Create(&newPartition).Error; err != nil {
		res.HttpStatus = http.StatusBadRequest
		res.Code = response.FailCode
		res.Msg = "创建分区失败"
		return res
	}

	return res
}

/*********************************************************
** 函数功能: 删除分区
** 日    期: 2021年12月9日
**********************************************************/
func DeletePartitionService(id uint) response.ResponseStruct {
	DB := common.GetDB()
	DB.Where("id = ?", id).Delete(model.Partition{})

	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       nil,
		Msg:        "ok",
	}
}

/*********************************************************
** 函数功能: 获取所有分区
** 日    期: 2021年12月9日
**********************************************************/
func GetAllPartitionService() response.ResponseStruct {
	var partitions []vo.AllPartitionVo

	DB := common.GetDB()
	DB.Model(&model.Partition{}).Select("id,content").Where("fid = 0").Scan(&partitions)
	for i := 0; i < len(partitions); i++ {
		//查询回复
		DB.Model(&model.Partition{}).Select("id,content").Where("fid = ?", partitions[i].ID).Scan(&partitions[i].Subpartition)
	}
	return response.ResponseStruct{
		HttpStatus: http.StatusOK,
		Code:       response.SuccessCode,
		Data:       gin.H{"partitions": partitions},
		Msg:        "ok",
	}
}

/*********************************************************
** 函数功能: 所属分区是否存在
** 日    期: 2021年12月9日
**********************************************************/
func IsParentPartitionExist(db *gorm.DB, fid uint) bool {
	var partition model.Partition
	db.Where("fid = 0").First(&partition, fid)
	if partition.ID != 0 {
		return true
	}
	return false
}

/*********************************************************
** 函数功能: 获取分区名
** 日    期: 2021年12月9日
**********************************************************/
func GetPartitionName(db *gorm.DB, id uint) string {
	var partition model.Partition
	var subpartition model.Partition
	db.First(&subpartition, id)
	if subpartition.ID != 0 {
		db.First(&partition, subpartition.Fid)
		return partition.Content + "/" + subpartition.Content
	}
	return "未分区"
}
