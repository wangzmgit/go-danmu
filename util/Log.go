package util

import (
	"io"
	"os"
	"time"
)

const (
	InfoLog  = "[Info] "
	ErrorLog = "[Error] "
)

//写入文件
func Logfile(logType string, log string) {
	var f1 *os.File
	var err error

	filename := "./file/logs/" + time.Now().Format("20060102") + ".log" //也可将name作为参数传进来

	if exist, _ := PathExists(filename); exist { //如果文件存在
		f1, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666) //打开文件,第二个参数是写入方式和权限
		if err != nil {
			PrintLog(ErrorLog, " Log file opening failed "+err.Error())
			return
		}
	} else {
		f1, err = os.Create(filename) //创建文件
		if err != nil {
			PrintLog(ErrorLog, " Log file creation failed "+err.Error())
			return
		}
	}
	_, err = io.WriteString(f1, logType+time.Now().Format("2006-01-02 15:04:05")+log+"\n") //写入文件
	if err != nil {
		PrintLog(ErrorLog, " Write log error "+err.Error())
	}
}

//打印日志
func PrintLog(logType string, log string) {
	println(logType + time.Now().Format("2006-01-02 15:04:05") + log)
}
