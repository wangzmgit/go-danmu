package model

import "github.com/jinzhu/gorm"

type Opinion struct {
	gorm.Model
	Name      string `gorm:"type:varchar(10);"`
	Email     string `gorm:"varchar(20);"`
	Telephone string `gorm:"varchar(20);"`
	Gender    int    `gorm:"default:0"`
	Desc      string `gorm:"varchar(200);"`
	Uid       uint
	Status    int `gorm:"default:0"` //状态
}
