package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-bluestore/log"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
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

func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		log.Error("sh -c eval echo ~$USER error.")
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		log.Error("blank output when reading home directory")
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		log.Error("Get home path error.")
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}

// Home returns the home directory for the executing user.
//
// This uses an OS-specific method for discovering the home directory.
// An error is returned if a home directory cannot be detected.
func Home() (string, error) {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}

	if "windows" == runtime.GOOS {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

// To ensure folder exist
func EnsureFolderExist(folderPath string) {
	exist, err := PathExists(folderPath)
	if nil != err {
		errMsg := fmt.Sprintf("access folder %s failed with %v", folderPath, err)
		panic(errMsg)
	}

	if !exist {
		err := os.Mkdir(folderPath, 0755)
		if err != nil {
			errMsg := fmt.Sprintf("create folder %s failed with %v", folderPath, err)
			panic(errMsg)
		}
	}
}
