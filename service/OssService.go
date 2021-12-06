package service

import (
	"container/list"

	"wzm/danmu3.0/util"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
)

type Callback func(vid int)

func UploadOSS(localFileName string, objectName string) (bool, string) {
	//储存到阿里云OSS
	client, err := oss.New(viper.GetString("aliyunoss.endpoint"), viper.GetString("aliyunoss.accesskey_id"), viper.GetString("aliyunoss.accesskey_secret"))
	if err != nil {
		util.Logfile("[Error]", " OSS请求错误 "+err.Error())
		return false, ""
	}
	// 获取存储空间
	bucket, err := client.Bucket(viper.GetString("aliyunoss.bucket"))
	if err != nil {
		util.Logfile("[Error]", " OSS请求错误 "+err.Error())
		return false, ""
	}

	err = bucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		util.Logfile("[Error]", " OSS上传失败 "+err.Error())
		return false, ""
	}
	//如果自定义域名长度为0使用默认地址
	var url string
	if len(viper.GetString("aliyunoss.domain")) == 0 {
		url = "http://" + viper.GetString("aliyunoss.bucket") + "." + viper.GetString("aliyunoss.endpoint") + "/" + objectName
	} else {
		url = "http://" + viper.GetString("aliyunoss.domain") + "/" + objectName
	}
	return true, url
}

func UploadVideoToOSS(localFileName string, objectName string, vid int) {
	//储存到阿里云OSS
	client, err := oss.New(viper.GetString("aliyunoss.endpoint"), viper.GetString("aliyunoss.accesskey_id"), viper.GetString("aliyunoss.accesskey_secret"))
	if err != nil {
		util.Logfile("[Error]", " OSS请求错误 "+err.Error())
		return
	}
	// 获取存储空间
	bucket, err := client.Bucket(viper.GetString("aliyunoss.bucket"))
	if err != nil {
		util.Logfile("[Error]", " OSS请求错误 "+err.Error())
		return
	}

	err = bucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		util.Logfile("[Error]", " OSS上传失败 "+err.Error())
		return
	}
	//完成上传
	CompleteUpload(vid)
}

func UploadFolderToOSS(dir string, files *list.List) bool {
	//储存到阿里云OSS
	client, err := oss.New(viper.GetString("aliyunoss.endpoint"), viper.GetString("aliyunoss.accesskey_id"), viper.GetString("aliyunoss.accesskey_secret"))
	if err != nil {
		return false
	}
	// 获取存储空间
	bucket, err := client.Bucket(viper.GetString("aliyunoss.bucket"))
	if err != nil {
		return false
	}
	objectName := "video/" + dir + "/"
	//上传m3u8文件
	err = bucket.PutObjectFromFile(objectName+"index.m3u8", "./file/output/"+dir+"/index.m3u8")
	if err != nil {
		return false
	}
	//上传ts文件
	for ts := files.Front(); ts != nil; ts = ts.Next() {
		oss := objectName + ts.Value.(string)
		local := "./file/output/" + dir + "/" + ts.Value.(string)
		err = bucket.PutObjectFromFile(oss, local)
		if err != nil {
			return false
		}
	}
	return true
}
