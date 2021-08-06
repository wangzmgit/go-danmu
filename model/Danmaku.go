package model

import "github.com/jinzhu/gorm"

type Danmaku struct {
	gorm.Model
	Vid     uint   `gorm:"not null;index"`
	Time    string `gorm:"not null"`  //时间
	Type    int    `gorm:"default:0"` //类型0滚动;1顶部;2底部
	Color   string `gorm:"type:varchar(10);default:'#fff'"`
	Content string `gorm:"type:varchar(100);not null"`
	Uid     uint   `gorm:"not null"`
}
