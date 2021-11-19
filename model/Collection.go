package model

type Collection struct {
	Cover  string `gorm:"size:255;not null"`
	Title  string `gorm:"type:varchar(50);not null;"`
	Desc   string `gorm:"size:255;"`
	Videos []Video
	Uid    uint // 作者
}
