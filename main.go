package main

import (
	"fmt"
	"github.com/go-bluestore/conf"
	"github.com/go-bluestore/log"
	"github.com/go-bluestore/utils"
	"os"
	"strings"
)

func InitLog(config log.Config) {

	if config.Appenders[conf.FileLogAppender].Enabled {
		// initialize logfile
		logPath := config.Appenders[conf.FileLogAppender].LogPath
		utils.EnsureFolderExist(logPath[0:strings.LastIndex(logPath, "/")])
		logfile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
		config.Appenders[conf.FileLogAppender].Output = logfile
	}

	log.SetGlobalConfig(&config)
}

func main() {
	config := conf.NewBlueStoreConfig()
	fmt.Printf("pares config %v.", config)

	InitLog(config.Logger)
	log.Debug("init log bluestore success.")
}
