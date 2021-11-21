package model

import "github.com/jinzhu/gorm"

type Collection struct {
	gorm.Model
	Cover string `gorm:"size:255;not null"`
	Title string `gorm:"type:varchar(50);not null;"`
	Desc  string `gorm:"size:255;"`
	Uid   uint   // 作者
}
