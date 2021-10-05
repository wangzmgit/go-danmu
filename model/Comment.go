package model

import (
	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model
	Vid     uint   `gorm:"not null;index"`             //视频ID
	Content string `gorm:"type:varchar(255);not null"` //内容
	Uid     uint   `gorm:"not null"`                   //用户
	//回复数,3.2版本新增,用于获取评论列表V2接口
	ReplyCount int `gorm:"default:0"`
}
