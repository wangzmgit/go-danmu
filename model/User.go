package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Avatar   string    `gorm:"size:255;"`
	Name     string    `gorm:"type:varchar(20);not null"`
	Email    string    `gorm:"varchar(20);not null;index"`
	Password string    `gorm:"size:255;not null"`
	Gender   int       `gorm:"default:0"`
	Birthday time.Time `gorm:"default:'1970-01-01'"`
	Sign     string    `gorm:"varchar(50);default:'这个人很懒，什么都没有留下'"`
}
