package service

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/spf13/viper"
	"kuukaa.fun/danmu-v4/util"
)

///转码hls
func Transcoding(completeName string, vid int) {
	var fileName, inputDir, outputDir, url string
	fileName = strings.TrimSuffix(completeName, path.Ext(completeName)) //无后缀文件名
	os.Mkdir("./file/output/"+fileName, os.ModePerm)

	//获取url
	if viper.GetBool("aliyunoss.storage") {
		url = GetUrl() + "video/" + fileName + "/"
	} else {
		url = GetUrl() + "output/" + fileName + "/"
	}
	//判断当前系统
	sysType := runtime.GOOS
	if sysType == "linux" {
		// LINUX系统
		inputDir = "./file/video/"
		outputDir = "./file/output/" + fileName + "/"
	} else if sysType == "windows" {
		// windows系统
		currentPath, _ := os.Getwd()
		inputDir = currentPath + "\\file\\video\\"
		outputDir = currentPath + "\\file\\output\\" + fileName + "\\"
	}
	//转到hls
	ToBitRate(inputDir+completeName, outputDir)
	ToHls(url, outputDir+"temp.ts", outputDir, fileName, vid)
}

/*********************************************************
** 函数功能: 降低码率转到ts
** 日    期: 2022年1月5日12:54:38
** 参    数: 视频宽度
**********************************************************/
func ToBitRate(inputFile string, outputDir string) error {
	outputFile := outputDir + "temp.ts"
	cmd := exec.Command("ffmpeg", "-y", "-i", inputFile,
		"-vcodec", "copy", "-acodec", "copy",
		"-vbsf", "h264_mp4toannexb", "-crf",
		"20", outputFile,
	)
	// 执行命令，返回命令是否执行成功
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}
	return nil
}

/*********************************************************
** 函数功能: 转到hls
** 日    期:
** 参    数: 文件域名,输入文件,输出目录,视频ID
**********************************************************/
func ToHls(url string, inputFile string, outputDir string, fileName string, vid int) {
	output := outputDir + "temp.m3u8"
	outputTs := outputDir + "output%03d.ts"
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-c",
		"copy", "-map", "0", "-f", "segment", "-segment_list",
		output, "-segment_time", "30", outputTs,
	)
	//fmt.Println(cmd.Args) //查看当前执行命令
	// 执行命令，返回命令是否执行成功
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	RewriteM3U8(url, output, fileName, vid)
}

//重写m3u8文件
func RewriteM3U8(url string, output string, fileName string, vid int) {
	file, err := os.OpenFile(output, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("open file filed.", err)
		return
	}
	//创建目标文件
	filename := "./file/output/" + fileName + "/index.m3u8" //也可将name作为参数传进来
	newFile, _ := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	//文件列表
	fileList := list.New()
	//defer关闭文件
	defer file.Close()
	defer newFile.Close()
	//读取文件内容到io中
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			//读到末尾
			if err == io.EOF {
				break
			} else {
				util.Logfile("[Error] ", "transcoding err "+err.Error())
				return
			}
		}
		//根据关键词覆盖当前行
		if strings.Contains(line, ".ts") {
			fileList.PushBack(strings.Replace(line, "\n", "", -1))
			newFile.WriteString(url + line)
		} else {
			newFile.WriteString(line)
		}
	}

	if viper.GetBool("aliyunoss.storage") {
		//将文件上传到oss
		success := UploadFolderToOSS(fileName, fileList)
		if !success {
			//上传失败，调用未通过审核
			VideoReviewFail(vid, "视频处理失败")
			return
		}
	}
	CompleteUpload(vid)
}
