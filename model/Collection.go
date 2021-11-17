package model

type Collection struct {
	Cover  string `gorm:"size:255;not null"`
	Title  string `gorm:"type:varchar(50);not null;"`
	Videos []Video
}
