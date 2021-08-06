package common

import (
	"fmt"
	"wzm/danmu3.0/model"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	driverName := viper.GetString("datasource.driverName")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)
	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to connect database ,err:" + err.Error())
	}
	//数据库迁移
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Video{})
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
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
