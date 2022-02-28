package util

import "os"

/*********************************************************
** 函数功能: 判断文件夹或文件是否存在
** 参    数：文件夹或文件路径
** 日    期: 2022年2月28日18:07:30
**********************************************************/
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
