package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-bluestore/log"
	"github.com/go-bluestore/utils"
)

const (
	// config file prefix
	ConfigPrefix = "bluestore"

	// Log Setting
	LogTimeFieldFormat = "logging.timeFieldFormat"
	ConsoleLogAppender = "logging.console"
	LogConsoleEnabled  = "logging.console.enabled"
	LogConsoleLevel    = "logging.console.level"
	LogConsoleFormat   = "logging.console.format"
	LogConsoleCaller   = "logging.console.caller"
	LogConsoleHostname = "logging.console.hostname"
	FileLogAppender    = "logging.file"
	LogFileEnabled     = "logging.file.enabled"
	LogFilePath        = "logging.file.path"
	LogFileLevel       = "logging.file.level"
	LogFileFormat      = "logging.file.format"
	LogFileCaller      = "logging.file.caller"
	LogFileHostname    = "logging.file.hostname"
)

type Config struct {
	filePath string
	maps     map[string]interface{}
}

type BlueStoreConfig struct {
	// log setting
	Logger log.Config
}

func LoadConfig() (config *viper.Viper) {
	config = viper.New()
	// for environment variables
	config.SetEnvPrefix(ConfigPrefix)
	config.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	config.SetEnvKeyReplacer(replacer)

	config.SetConfigName(ConfigPrefix)
	homePath, _ := utils.Home()
	config.AddConfigPath(fmt.Sprintf("%s/.%s", homePath, ConfigPrefix))
	// Path to look for the config file in based on GOPATHc
	goPath := os.Getenv("GOPATH")
	for _, p := range filepath.SplitList(goPath) {
		config.AddConfigPath(filepath.Join(p, "src/github.com/go-bluestore"))
	}

	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error reading plugin config: %s", err))
	}
	return
}

func NewBlueStoreConfig() BlueStoreConfig {
	config := LoadConfig()

	return BlueStoreConfig{
		Logger:GetLogSetting(config),
	}
}

func GetLogSetting(conf *viper.Viper) log.Config {
	logTimestampFormat := conf.GetString(LogTimeFieldFormat)
	logConsoleEnabled := conf.GetBool(LogConsoleEnabled)
	logConsoleLevel := conf.GetInt(LogConsoleLevel)
	logConsoleFormat := conf.GetString(LogConsoleFormat)
	logConsoleCaller := conf.GetBool(LogConsoleCaller)
	logConsoleHostname := conf.GetBool(LogConsoleHostname)
	logFileEnabled := conf.GetBool(LogFileEnabled)
	logFilePath := conf.GetString(LogFilePath)
	logFileLevel := conf.GetInt(LogFileLevel)
	logFileFormat := conf.GetString(LogFileFormat)
	logFileCaller := conf.GetBool(LogFileCaller)
	logFileHostname := conf.GetBool(LogFileHostname)

	consoleAppender := &log.Appender{
		Enabled:      logConsoleEnabled,
		LogLevel:     log.Level(logConsoleLevel),
		LogType:      log.ConsoleLog,
		LogPath:      log.ConsoleStdout,
		Output:       os.Stdout,
		Format:       strings.ToUpper(logConsoleFormat),
		ShowCaller:   logConsoleCaller,
		ShowHostname: logConsoleHostname,
	}
	//tools.EnsureFolderExist(logFilePath[0:strings.LastIndex(logFilePath, "/")])
	//logfile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	//if err != nil {
	//	panic(err)
	//}
	fileAppender := &log.Appender{
		Enabled:      logFileEnabled,
		LogLevel:     log.Level(logFileLevel),
		LogType:      log.FileLog,
		LogPath:      logFilePath,
		Output:       nil,
		Format:       strings.ToUpper(logFileFormat),
		ShowCaller:   logFileCaller,
		ShowHostname: logFileHostname,
	}

	logConfig := log.Config{
		Enabled:         logConsoleEnabled || logFileEnabled,
		Provider:        log.GetGlobalConfig().Provider,
		GlobalLogLevel:  log.Level(uint8(math.Max(float64(logConsoleLevel), float64(logFileLevel)))),
		TimeFieldFormat: logTimestampFormat,
		Appenders:       map[string]*log.Appender{ConsoleLogAppender: consoleAppender, FileLogAppender: fileAppender},
		OutputFlags:     log.GetOutputFlags(),
	}
	return logConfig
}

