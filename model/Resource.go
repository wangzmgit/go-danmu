package model

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

/*视频资源表，
 *一个视频可能有多个分P
 *一个分P也有多个分辨率
 */
type Resource struct {
	gorm.Model
	//资源标识
	UUID uuid.UUID `gorm:"char(36);not null;"`
	//所属视频
	Vid uint
	//分P使用的标题
	Title string `gorm:"varchar(20);"`
	//不同分辨率
	Res360  string `gorm:"varchar(255);"`
	Res480  string `gorm:"varchar(255);"`
	Res720  string `gorm:"varchar(255);"`
	Res1080 string `gorm:"varchar(255);"`
	//原始分辨率，适用于早期版本未指定分辨率的视频
	//或者不进行转码处理的情况
	Original string `gorm:"varchar(255);"`
}

func (resource *Resource) BeforeCreate(tx *gorm.DB) (err error) {
	resource.UUID = uuid.New()
	return
}
