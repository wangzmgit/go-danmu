package util

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
)

func UploadOSS(localFileName string, objectName string) (bool, string) {
	//储存到阿里云OSS
	client, err := oss.New(viper.GetString("aliyunoss.endpoint"), viper.GetString("aliyunoss.accessid"), viper.GetString("aliyunoss.accesskey"))
	if err != nil {
		return false, "OSS请求错误" + err.Error()
	}
	// 获取存储空间
	bucket, err := client.Bucket(viper.GetString("aliyunoss.bucket"))
	if err != nil {
		return false, "OSS请求错误" + err.Error()
	}

	err = bucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		return false, "OSS上传失败" + err.Error()
	}
	url := "http://" + viper.GetString("aliyunoss.bucket") + "." + viper.GetString("aliyunoss.endpoint") + "/" + objectName
	return true, url
}
