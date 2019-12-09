package conf

import (
	"fmt"
	"github.com/spf13/viper"
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
		Logger: GetLogSetting(config),
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
		Enabled:         true,
		Provider:        log.GetGlobalConfig().Provider,
		TimeFieldFormat: logTimestampFormat,
	}

	if logConsoleEnabled {
		if logFileEnabled {
			panic("can not set console and file at the same time. ")
		}

		logConfig.OutputFlags = log.GetOutputFlags()
		logConfig.GlobalLogLevel = log.Level(uint8(float64(logConsoleLevel)))
		logConfig.Appenders = map[string]*log.Appender{ConsoleLogAppender: consoleAppender}
	} else {
		if !logFileEnabled {
			panic("choose from one of console or file at least. ")
		}

		logConfig.OutputFlags = log.GetOutputFlags()
		logConfig.GlobalLogLevel = log.Level(uint8(float64(logFileLevel)))
		logConfig.Appenders = map[string]*log.Appender{FileLogAppender: fileAppender}
	}

	return logConfig
}
