package log

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Logger interface defines all behaviors of a backendLogger.
type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
	Panic(msg string)

	DebugKV(msg string, keyvals map[string]interface{})
	InfoKV(msg string, keyvals map[string]interface{})
	WarnKV(msg string, keyvals map[string]interface{})
	ErrorKV(msg string, keyvals map[string]interface{})
	FatalKV(msg string, keyvals map[string]interface{})
	PanicKV(msg string, keyvals map[string]interface{})

	SetGlobalLogLevel(level Level)
	SetOutputFlags(flags *OutputFlags)
	SetTimeFieldFormat(format string)
}

// backendLogger is the actual logging object that we use, and is pre-built as a global variable.
var backendLogger = globalConfig.buildLogger()

// Disable to stop logging.
func Disable() {
	globalConfig.Enabled = false
	backendLogger = &nonLogger{}
}

// Enable to start logging.
func Enable() {
	globalConfig.Enabled = true
	backendLogger = globalConfig.buildLogger()
}

func Debug(fmtmsg string, a ...interface{}) {
	backendLogger.Debug(fmt.Sprintf(fmtmsg, a...))
}

func Info(fmtmsg string, a ...interface{}) {
	backendLogger.Info(fmt.Sprintf(fmtmsg, a...))
}

func Warn(fmtmsg string, a ...interface{}) {
	backendLogger.Warn(fmt.Sprintf(fmtmsg, a...))
}

func Error(fmtmsg string, a ...interface{}) {
	backendLogger.Error(fmt.Sprintf(fmtmsg, a...))
}

func Fatal(fmtmsg string, a ...interface{}) {
	backendLogger.Fatal(fmt.Sprintf(fmtmsg, a...))
}

func Panic(fmtmsg string, a ...interface{}) {
	backendLogger.Panic(fmt.Sprintf(fmtmsg, a...))
}

func DebugKV(msg string, keyvals map[string]interface{}) {
	backendLogger.DebugKV(msg, keyvals)
}

func InfoKV(msg string, keyvals map[string]interface{}) {
	backendLogger.InfoKV(msg, keyvals)
}

func WarnKV(msg string, keyvals map[string]interface{}) {
	backendLogger.WarnKV(msg, keyvals)
}

func ErrorKV(msg string, keyvals map[string]interface{}) {
	backendLogger.ErrorKV(msg, keyvals)
}

func FatalKV(msg string, keyvals map[string]interface{}) {
	backendLogger.FatalKV(msg, keyvals)
}

func PanicKV(msg string, keyvals map[string]interface{}) {
	backendLogger.PanicKV(msg, keyvals)
}

func GetGlobalConfig() *Config {
	return globalConfig
}

// SetGlobalConfig is used to refresh logging manners.
func SetGlobalConfig(config *Config) {
	globalConfig = config
	backendLogger = globalConfig.buildLogger()
}

func GetGlobalLogLevel() Level {
	return globalConfig.GlobalLogLevel
}

// SetGlobalLogLevel is used to restraint log-level of all "Appenders".
func SetGlobalLogLevel(level Level) {
	globalConfig.GlobalLogLevel = level
	backendLogger.SetGlobalLogLevel(level)
}

func GetOutputFlags() *OutputFlags {
	return globalConfig.OutputFlags
}

// SetOutputFlags is used to reconfig output flags.
func SetOutputFlags(flags *OutputFlags) {
	globalConfig.OutputFlags = flags
	backendLogger.SetOutputFlags(flags)
}

func SetTimestampFormat(format string) {
	globalConfig.TimeFieldFormat = format
	backendLogger.SetTimeFieldFormat(format)
}

// AddAppender adds/replaces a new logging destination.
func AddAppender(appenderName string, output io.Writer, logLevel Level, format string, showCaller bool, showHostname bool) {
	globalConfig.Appenders[appenderName] = &Appender{
		Enabled:      true,
		LogLevel:     logLevel,
		Output:       output,
		Format:       format,
		ShowCaller:   showCaller,
		ShowHostname: showHostname,
	}
	backendLogger = globalConfig.buildLogger()
}

// RemoveAppender removes a logging appender by name.
func RemoveAppender(appenderNameToRemove string) {
	delete(globalConfig.Appenders, appenderNameToRemove)
	backendLogger = globalConfig.buildLogger()
}

// AddFileAppender adds/replaces a new logging destination that append logs to a specified file.
func AddFileAppender(appenderName string, filePath string, logLevel Level, format string, showCaller bool, showHostname bool) {
	_, err := os.Stat(filePath)

	if err != nil {
		if os.IsNotExist(err) {
			parentPath := filePath[0:strings.LastIndex(filePath, "/")]
			_, err := os.Stat(filePath)
			if err != nil {
				if os.IsNotExist(err) {
					err := os.MkdirAll(parentPath, 0755)
					if err != nil {
						panic(err)
					}
				}
			}
		} else {
			panic(err)
		}
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	AddAppender(appenderName, file, logLevel, format, showCaller, showHostname)
}

// SetAppenders sets a set of "Appenders".
func SetAppenders(appenders map[string]*Appender) {
	globalConfig.Appenders = appenders
	backendLogger = globalConfig.buildLogger()
}
