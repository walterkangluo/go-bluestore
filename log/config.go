package log

import (
	"io"
	"os"
)

// Appender is responsible for delivering LogEvents to their destination.
type Appender struct {
	Enabled      bool
	LogLevel     Level
	LogType      string
	LogPath      string
	Output       io.Writer
	Format       string
	ShowCaller   bool
	ShowHostname bool
}

// OutputFlags are printed in log record that can be customized.
type OutputFlags struct {
	// TimestampFieldName is the field name used for the timestamp field.
	TimestampFieldName string
	// LevelFieldName is the field name used for the level field.
	LevelFieldName string
	// MessageFieldName is the field name used for the message field.
	MessageFieldName string
	// ErrorFieldName is the field name used for error fields.
	ErrorFieldName string
	// CallerFieldName is the field name used for caller field.
	CallerFieldName string
	// HostnameFieldName is the field name used for hostname field.
	HostnameFieldName string
}

// Config includes configurations for our log, such as log-level.
// For more log destinations just add "Appender" into "Config.[]Appenders".
type Config struct {
	Enabled         bool
	Provider        Provider
	GlobalLogLevel  Level
	TimeFieldFormat string
	Appenders       map[string]*Appender
	OutputFlags     *OutputFlags
}

// stdoutAppender is a pre-configed console log.
var stdoutAppender = &Appender{
	Enabled:      true,
	LogLevel:     InfoLevel,
	LogType:      ConsoleLog,
	LogPath:      ConsoleStdout,
	Output:       os.Stdout,
	Format:       TextFmt,
	ShowCaller:   true,
	ShowHostname: true,
}

// globalOutputFlags contains pre-defined output flags. Usually no need to modify.
var globalOutputFlags = &OutputFlags{
	TimestampFieldName: "time",
	LevelFieldName:     "level",
	MessageFieldName:   "message",
	ErrorFieldName:     "error",
	CallerFieldName:    "caller",
	HostnameFieldName:  "host",
}

// globalConfig is a set of default log configuration with only one "stdoutAppender".
var globalConfig = &Config{
	Enabled:         true,
	Provider:        Zerolog,
	GlobalLogLevel:  DebugLevel,
	TimeFieldFormat: "2006-01-02 15:04:05.000",
	Appenders:       map[string]*Appender{"stdout": stdoutAppender},
	OutputFlags:     globalOutputFlags,
}

// buildLogger builds a "Logger" with a number of backend logger inside.
// Each logger corresponds to an "Appender".
func (config *Config) buildLogger() Logger {
	if !config.Enabled {
		return nonLogger{}
	}
	switch config.Provider {
	case Zerolog:
		logger := buildZeroLogger(config)
		return logger
	}
	return nil
}

// Level defines log levels.
type Level uint8

const (
	// DebugLevel defines debug log level.
	DebugLevel Level = iota
	// InfoLevel defines info log level.
	InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel
	// FatalLevel defines fatal log level.
	FatalLevel
	// PanicLevel defines panic log level.
	PanicLevel
	// Disabled disables the logger.
	Disabled
)

// Provider enumerates backend log libs.
type Provider uint8

const (
	Zerolog Provider = iota
)

const (
	// JsonFmt indicates that log output generated in form of JSON.
	JsonFmt string = "JSON"
	// TextFmt indicates that log output generated in form of TEXT.
	TextFmt string = "TEXT"
	// ConsoleLog indicates that log output to console.
	ConsoleLog string = "CONSOLE_LOG"
	// FileLog indicates that log output to console.
	FileLog string = "FILE_LOG"
	// ConsoleStdout indicates than console log output to os.Stdout
	ConsoleStdout string = "STDOUT"
	// ConsoleStderr indicates than console log output to os.Stderr
	ConsoleStderr string = "STDERR"
)
