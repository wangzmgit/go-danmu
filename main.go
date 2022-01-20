package main

import (
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"kuukaa.fun/danmu-v4/common"
	"kuukaa.fun/danmu-v4/routes"
	"kuukaa.fun/danmu-v4/util"
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
	// 创建gin日志文件
	file := InitGinLog()
	gin.DisableConsoleColor()
	gin.DefaultWriter = io.MultiWriter(file)
	// 设置模式
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

func InitGinLog() *os.File {
	filenames := "./file/logs/gin_" + util.RandomString(3) + time.Now().Format("20060102") + ".log"
	file, err := os.Create(filenames)
	if err != nil {
		return nil
	}
	return file
}
