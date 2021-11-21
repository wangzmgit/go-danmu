package model

import "github.com/jinzhu/gorm"

type VideoCollection struct {
	gorm.Model
	Vid          uint
	CollectionId uint
}
