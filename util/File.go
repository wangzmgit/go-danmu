package util

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

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

/*********************************************************
** 函数功能: 获取指定目录下的子目录
** 参    数：目录的路径
** 日    期: 2022年3月10日11:06:39
**********************************************************/
func GetSubDir(path string) ([]string, error) {
	dirs := make([]string, 0)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		Logfile(ErrorLog, "read skins error: "+err.Error())
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		}
	}
	return dirs, nil
}

/*********************************************************
** 函数功能: 解压zip文件
** 参    数：输入、输出路径
** 日    期: 2022年3月10日13:00:28
**********************************************************/
func Unzip(zipFile string, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
