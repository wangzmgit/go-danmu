package model

import "github.com/jinzhu/gorm"

type Partition struct {
	gorm.Model
	Content string `gorm:"varchar(20);"`
	Fid     uint   //所属分区ID
}
