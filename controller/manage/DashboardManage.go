package manage

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"kuukaa.fun/danmu-v4/common"
	"kuukaa.fun/danmu-v4/response"
	"kuukaa.fun/danmu-v4/service"
	"kuukaa.fun/danmu-v4/util"
)

/*********************************************************
** 函数功能: 获取网站数据
** 日    期: 2021年11月25日18:00:40
**********************************************************/
func GetTotalWebsiteData(ctx *gin.Context) {
	res := service.GetTotalWebsiteDataService()
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 获取网站近期(5天)数据
** 日    期: 2021年11月25日
**********************************************************/
func GetRecentWebsiteData(ctx *gin.Context) {
	res := service.GetGetRecentWebsiteDataService()
	response.HandleResponse(ctx, res)
}

/*********************************************************
** 函数功能: 检查更新
** 日    期: 2022年3月11日12:19:31
**********************************************************/
func CheckUpdate(ctx *gin.Context) {
	//请求仓库地址
	url := "https://gitee.com/api/v5/repos/wzmgit/go-danmu/contents/common%2FVersion.go"
	res, err := http.DefaultClient.Get(url)
	if err != nil {
		response.Fail(ctx, nil, response.CheckUpdateFail)
		return
	}
	defer res.Body.Close()

	var versionInfo UpdateVersion
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		response.Fail(ctx, nil, response.CheckUpdateFail)
		return
	}
	// 解析JSON
	err = json.Unmarshal(body, &versionInfo)
	if err != nil {
		response.Fail(ctx, nil, response.CheckUpdateFail)
		return
	}

	//解码base64
	byteInfo, err := base64.StdEncoding.DecodeString(versionInfo.Content)
	if err != nil {
		response.Fail(ctx, nil, response.RequestError)
		return
	}
	//正则匹配版本
	pattern := `([0-9]+).([0-9]+).([0-9]+)`
	reg := regexp.MustCompile(pattern)
	remoteVersion := reg.FindStringSubmatch(string(byteInfo)) //远程版本
	localVersion := reg.FindStringSubmatch(common.Version)    //本地版本
	if len(remoteVersion) < 4 || len(localVersion) < 4 {
		response.Fail(ctx, nil, response.RequestError)
		return
	}
	//比对版本信息
	for i := 1; i < 4; i++ {
		if util.StringCompare(remoteVersion[i], localVersion[i]) {
			response.Success(ctx, gin.H{"version": remoteVersion[0]}, response.OK)
			return
		}
	}
	response.Success(ctx, gin.H{"version": nil}, response.OK)
}
