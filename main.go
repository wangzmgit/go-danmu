package main

import (
	"os"
	"wzm/danmu3.0/common"
	"wzm/danmu3.0/cronJob"
	"wzm/danmu3.0/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

const (
	// DebugMode indicates gin mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates gin mode is release.
	ReleaseMode = "release"
	// TestMode indicates gin mode is test.
	TestMode = "test"
)

func main() {
	InitConfig()
	println("     _")
	println("  __| | __ _ _ __  _ __ ___  _   _")
	println(" / _` |/ _` | '_ \\| '_ ` _ \\| | | |")
	println("| (_| | (_| | | | | | | | | | |_| |")
	print(" \\__,_|\\__,_|_| |_|_| |_| |_|\\__,_|")
	println("\tversion:" + common.Version)
	//初始化Redis
	common.Redis()
	//初始化数据库
	db := common.InitDB()
	defer db.Close()
	//创建定时任务
	cronJob.CronJob()
	gin.SetMode(ReleaseMode)
	r := gin.Default()
	r = routes.CollectRoute(r)
	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}
	panic(r.Run())
}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
