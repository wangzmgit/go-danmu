package model

import "github.com/jinzhu/gorm"

type Partition struct {
	gorm.Model
	Content string `gorm:"varchar(20);"`
	Fid     uint   `gorm:"default:0"` //所属分区ID
}
