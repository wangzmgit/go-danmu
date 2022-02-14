package service

import (
	"bufio"
	"bytes"
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/viper"
	"kuukaa.fun/danmu-v4/dto"
	"kuukaa.fun/danmu-v4/util"
)

var resInfo = map[int]string{
	360:  "600:360",
	480:  "720:480",
	720:  "1280:720",
	1080: "1920:1080",
}

var resName = map[int]string{
	0:    "original",
	360:  "360p",
	480:  "480p",
	720:  "720p",
	1080: "1080p",
}

///转码hls
func Transcoding(completeName string, vid int, maxRes int) {
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

	if maxRes == 0 {
		//调整码率
		ToBitRate(inputDir+completeName, outputDir, 0)
		//只转码原始分辨率
		ToHls(url, outputDir+"temp_"+resName[0]+".ts", outputDir, fileName, vid)
	} else {
		var wg sync.WaitGroup
		CreateResDir(maxRes, fileName) //创建不同分辨率文件夹
		switch maxRes {
		case 1080:
			wg.Add(1)
			go func() {
				file1080, _ := ToTargetRes(inputDir+completeName, outputDir, 1080)
				ToBitRate(file1080, outputDir, 1080)
				ToHlsDifferentRes(url, outputDir+"temp_"+resName[1080]+".ts", outputDir, fileName, vid, 1080)
				wg.Done()
			}()
			fallthrough
		case 720:
			wg.Add(1)
			go func() {
				file720, _ := ToTargetRes(inputDir+completeName, outputDir, 720)
				ToBitRate(file720, outputDir, 720)
				ToHlsDifferentRes(url, outputDir+"temp_"+resName[720]+".ts", outputDir, fileName, vid, 720)
				wg.Done()
			}()
			fallthrough
		case 480:
			wg.Add(1)
			go func() {
				file480, _ := ToTargetRes(inputDir+completeName, outputDir, 480)
				ToBitRate(file480, outputDir, 480)
				ToHlsDifferentRes(url, outputDir+"temp_"+resName[480]+".ts", outputDir, fileName, vid, 480)
				wg.Done()
			}()
			fallthrough
		case 360:
			wg.Add(1)
			go func() {
				file360, _ := ToTargetRes(inputDir+completeName, outputDir, 360)
				ToBitRate(file360, outputDir, 360)
				ToHlsDifferentRes(url, outputDir+"temp_"+resName[360]+".ts", outputDir, fileName, vid, 360)
				wg.Done()
			}()
		}
		wg.Wait()
		CompleteUpload(vid)
	}
	DeleteTempFile(maxRes, fileName) //删除临时文件
}

/*********************************************************
** 函数功能: 预处理视频
** 日    期: 2022年2月13日11:06:35
** 参    数: 最大分辨率
**********************************************************/
func PreTreatmentVideo(input string, vid int) (int, error) {
	var err error
	var videoData dto.VideoInfoData
	videoData, err = GetVideoInfo(input)
	if err != nil {
		return 0, err
	}

	if videoData.Stream[0].CodecName != "h264" {
		return 0, err
	}
	//计算最大分辨率
	width := videoData.Stream[0].Width
	height := videoData.Stream[0].Height
	maxRes := util.Min(GetWidthRes(width), GetHeigthRes(height))

	return maxRes, nil
}

/*********************************************************
** 函数功能: 获取视频信息
** 日    期: 2022年1月4日17:23:40
**********************************************************/
func GetVideoInfo(input string) (dto.VideoInfoData, error) {
	var err error
	var videoData dto.VideoInfoData
	// input = "./file/video/" + input + ".mp4"
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", input)
	// 执行命令，返回命令是否执行成功
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误
	err = cmd.Run()
	if err != nil {
		return videoData, errors.New(fmt.Sprint(err) + ":" + stderr.String())
	}
	// fmt.Println("Result: " + stdout.String())

	err = json.Unmarshal(stdout.Bytes(), &videoData)
	if err != nil {
		return videoData, err
	}

	return videoData, nil
}

/*********************************************************
** 函数功能: 获取宽度支持的最大分辨率
** 日    期: 2022年1月5日9:28:36
** 参    数: 视频宽度
**********************************************************/
func GetWidthRes(width int) int {
	//1920*1080
	if width >= 1920 {
		return 1080
	}
	// 1280*720
	if width >= 1280 {
		return 720
	}
	//720*480
	if width >= 720 {
		return 480
	}
	return 360
}

/*********************************************************
** 函数功能: 获取高度支持的最大分辨率
** 日    期: 2022年1月5日9:30:05
** 参    数: 视频高度
**********************************************************/
func GetHeigthRes(height int) int {
	//1920*1080
	if height >= 1080 {
		return 1080
	}
	// 1280*720
	if height >= 720 {
		return 720
	}
	//720*480
	if height >= 480 {
		return 480
	}
	return 360
}

/*********************************************************
** 函数功能: 降低码率转到ts
** 日    期: 2022年1月5日12:54:38
** 参    数: 输入文件，输出目录，分辨率
**********************************************************/
func ToBitRate(inputFile string, outputDir string, res int) error {
	outputFile := outputDir + "temp_" + resName[res] + ".ts"
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

/*********************************************************
** 函数功能: 转到hls(指定分辨率)
** 日    期:
** 参    数: 文件域名,输入文件,输出目录,视频ID
**********************************************************/
func ToHlsDifferentRes(url string, inputFile string, outputDir string, fileName string, vid int, res int) {
	output := outputDir + strconv.Itoa(res) + "p/temp.m3u8"
	outputTs := outputDir + strconv.Itoa(res) + "p/output%03d.ts"
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-c", "copy",
		"-map", "0", "-f", "segment", "-segment_list",
		output, "-segment_time", "30", outputTs,
	)
	// fmt.Println(cmd.Args) //查看当前执行命令
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
	RewriteDifferentRes(url, output, fileName, vid, res)
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

/*********************************************************
** 函数功能: 重写m3u8文件,多分辨率
** 日    期: 2022年2月13日11:56:53
** 参    数: 文件域名,输入文件,输出目录,视频ID,分辨率
**********************************************************/
func RewriteDifferentRes(url string, output string, fileName string, vid int, res int) {
	newUrl := url + strconv.Itoa(res) + "p/"
	file, err := os.OpenFile(output, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("open file filed.", err)
		return
	}
	//创建目标文件
	filename := "./file/output/" + fileName + "/" + strconv.Itoa(res) + "p/index.m3u8"
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
			newFile.WriteString(newUrl + line)
		} else {
			newFile.WriteString(line)
		}
	}

	if viper.GetBool("aliyunoss.storage") {
		//将文件上传到oss
		success := UploadFolderToOSS(fileName, fileList)
		if !success {
			//上传失败，调用未通过审核
			VideoReviewFail(vid, "视频上传失败")
			return
		}
	}
}

/*********************************************************
** 函数功能: 转到目标分辨率
** 日    期: 2022年2月13日17:46:59
** 参    数: 输入文件,输出目录,分辨率
**********************************************************/
func ToTargetRes(inputFile string, outputDir string, res int) (string, error) {
	outputFile := outputDir + "temp_" + resName[res] + ".mp4"
	cmd := exec.Command("ffmpeg", "-i", inputFile,
		"-vf", "scale="+resInfo[res], outputFile,
		"-hide_banner",
	)
	// 执行命令，返回命令是否执行成功
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return "", err
	}
	return outputFile, nil
}
