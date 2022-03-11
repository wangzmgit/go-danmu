package common

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"kuukaa.fun/danmu-v4/model"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		username,
		password,
		host,
		port,
		database)
	db, err := gorm.Open("mysql", args)
	if err != nil {
		panic("failed to connect database ,err:" + err.Error())
	}
	//数据库迁移
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Video{})
	db.AutoMigrate(&model.Resource{})
	db.AutoMigrate(&model.Review{})
	db.AutoMigrate(&model.Interactive{})
	db.AutoMigrate(&model.Follow{})
	db.AutoMigrate(&model.Comment{})
	db.AutoMigrate(&model.Reply{})
	db.AutoMigrate(&model.Announce{})
	db.AutoMigrate(&model.AnnounceUser{})
	db.AutoMigrate(&model.Message{})
	db.AutoMigrate(&model.Danmaku{})
	db.AutoMigrate(&model.Carousel{})
	db.AutoMigrate(&model.Admin{})
	db.AutoMigrate(&model.Collection{})
	db.AutoMigrate(&model.VideoCollection{})
	db.AutoMigrate(&model.Opinion{})
	db.AutoMigrate(&model.Partition{})
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
