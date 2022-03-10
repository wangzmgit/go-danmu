package manage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/dto"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/util"
)

type SkinInfo struct {
	FileName string `json:"file_name"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	Cover    string `json:"cover"`
	Author   string `json:"author"`
	Desc     string `json:"desc"`
}

func GetSkinInfoList(ctx *gin.Context) {
	var skins []SkinInfo
	//读取主题文件夹下所有目录
	skinFiles, err := util.GetSubDir("./file/skins")
	if err != nil {
		response.Fail(ctx, nil, response.ReadFail)
		return
	}
	//读取目录下的配置文件
	for i := 0; i < len(skinFiles); i++ {
		infoFile, err := os.Open("./file/skins/" + skinFiles[i] + "/info.json")
		if err != nil {
			continue
		}
		defer infoFile.Close()

		byteValue, _ := ioutil.ReadAll(infoFile)
		var skinInfo SkinInfo
		err = json.Unmarshal(byteValue, &skinInfo)
		if err != nil {
			continue
		}
		skinInfo.FileName = skinFiles[i]
		skins = append(skins, skinInfo)
	}

	response.Success(ctx, gin.H{"skins": skins}, response.OK)
}

func ApplySkin(ctx *gin.Context) {
	var request dto.SkinDto
	requestErr := ctx.Bind(&request)
	if requestErr != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}
	//检测操作系统
	sysType := runtime.GOOS
	if sysType != "linux" {
		response.Fail(ctx, nil, response.SystemNotSupported)
		return
	}
	fileName := request.FileName + ".zip"
	// 文件是否存在
	exist, _ := util.PathExists("./file/skins/" + fileName)
	if !exist {
		response.Fail(ctx, nil, response.SkinNotExist)
		return
	}

	//移除旧版主题
	os.RemoveAll("/var/www/danmaku")
	//解压新的主题到目标目录
	util.Unzip("./file/skins/"+fileName, "/var/www/danmaku")
	response.Success(ctx, nil, response.OK)
}

func DeleteSkin(ctx *gin.Context) {
	var request dto.SkinDto
	requestErr := ctx.Bind(&request)
	if requestErr != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}
	//删除压缩文件
	os.Remove("./file/skins/" + request.FileName + ".zip")
	//删除解压后的文件
	os.RemoveAll("./file/skins/" + request.FileName)
	response.Success(ctx, nil, response.OK)
}
