package model

import "github.com/jinzhu/gorm"

type VideoResource struct {
	gorm.Model
	Vid        uint   `gorm:"not null;index"`    //所属视频
	Resolution string `gorm:"type:varchar(20);"` //分辨率信息
	url        string `gorm:"size:255;not null"` //资源
}
