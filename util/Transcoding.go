package util

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
)

///转码hls
func Transcoding(name string, vid int, callback Callback) string {
	var input, output, dir, url string
	dir = strings.TrimSuffix(name, path.Ext(name)) //文件夹名
	os.Mkdir("./file/output/"+dir, os.ModePerm)
	//生成url
	url = "http://" + viper.GetString("aliyunoss.bucket") + "." + viper.GetString("aliyunoss.endpoint") + "/video/" + dir + "/"
	//判断当前系统
	sysType := runtime.GOOS
	if sysType == "linux" {
		// LINUX系统
		input = "./file/video/" + name
		output = "./file/output/" + dir + "/" + dir + ".m3u8"
	} else if sysType == "windows" {
		// windows系统
		currentPath, _ := os.Getwd()
		input = currentPath + "\\file\\video\\" + name
		output = currentPath + "\\file\\output\\" + dir + "\\" + dir + ".m3u8"
	}
	//转到hls
	go ToHls(url, input, output, dir, vid, callback)
	return url + "index.m3u8"
}

func ToHls(url string, input string, output string, dir string, vid int, callback Callback) {
	cmd := exec.Command("ffmpeg", "-i", input, "-c:v", "libx264", "-c:a", "aac", "-strict", "-2", "-f", "hls", "-hls_list_size", "0", "-hls_time", "15", output)
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
	fmt.Println("Result: " + out.String())
	RewriteM3U8(url, output, dir, vid, callback)

}

func RewriteM3U8(url string, path string, dir string, vid int, callback Callback) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("open file filed.", err)
		return
	}
	//创建目标文件
	filename := "./file/output/" + dir + "/index.m3u8" //也可将name作为参数传进来
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
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
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
	//将文件上传到oss
	success := UploadFolderToOSS(dir, fileList)
	if success {
		//完成上传
		ProcessingComplete(vid, callback)
	}
}
