package utils

import (
	"os"
)

// 判断路径是否存在，不论是文件或者目录
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if nil != err {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// 判断输入的路径是否为目录
// 注意：这里不判断是否存在该目录
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if nil != err {
		return false
	}

	return s.IsDir()
}

// 判断输入路径是否为文件
// 注意：这里不判断是否存在该文件
func IsFile(path string) bool {
	return !IsDir(path)
}


