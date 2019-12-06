package log

import (
	"os"

	"github.com/rs/zerolog"
)

// Interface assertions
var _ Logger = (*zeroLogger)(nil)

// zeroLogger is used as an adapter which contains a set of loggers,
// each logger corresponds to an "Appender".
type zeroLogger struct {
	loggers map[string]*zerolog.Logger
}

func (zl *zeroLogger) SetOutputFlags(flags *OutputFlags) {
	zerolog.CallerFieldName = flags.CallerFieldName
	zerolog.ErrorFieldName = flags.ErrorFieldName
	zerolog.LevelFieldName = flags.LevelFieldName
	zerolog.MessageFieldName = flags.MessageFieldName
	zerolog.TimestampFieldName = flags.TimestampFieldName
}

func (zl *zeroLogger) SetGlobalLogLevel(level Level) {
	zerolog.SetGlobalLevel(parseLogLevel(level))
}

func (zl *zeroLogger) SetTimeFieldFormat(format string) {
	zerolog.TimeFieldFormat = format
}

func (zl *zeroLogger) Debug(msg string) {
	for _, l := range zl.loggers {
		l.Debug().Msg(msg)
	}
}

func (zl *zeroLogger) Info(msg string) {
	for _, l := range zl.loggers {
		l.Info().Msg(msg)
	}
}

func (zl *zeroLogger) Warn(msg string) {
	for _, l := range zl.loggers {
		l.Warn().Msg(msg)
	}
}

func (zl *zeroLogger) Error(msg string) {
	for _, l := range zl.loggers {
		l.Error().Msg(msg)
	}
}

func (zl *zeroLogger) Fatal(msg string) {
	for _, l := range zl.loggers {
		l.Fatal().Msg(msg)
	}
}

func (zl *zeroLogger) Panic(msg string) {
	for _, l := range zl.loggers {
		l.Panic().Msg(msg)
	}
}

func (zl *zeroLogger) DebugKV(msg string, keyvals map[string]interface{}) {
	for _, l := range zl.loggers {
		l.Debug().Fields(keyvals).Msg(msg)
	}
}

func (zl *zeroLogger) InfoKV(msg string, keyvals map[string]interface{}) {
	for _, l := range zl.loggers {
		l.Info().Fields(keyvals).Msg(msg)
	}
}

func (zl *zeroLogger) WarnKV(msg string, keyvals map[string]interface{}) {
	for _, l := range zl.loggers {
		l.Warn().Fields(keyvals).Msg(msg)
	}
}

func (zl *zeroLogger) ErrorKV(msg string, keyvals map[string]interface{}) {
	for _, l := range zl.loggers {
		l.Error().Fields(keyvals).Msg(msg)
	}
}

func (zl *zeroLogger) FatalKV(msg string, keyvals map[string]interface{}) {
	for _, l := range zl.loggers {
		l.Fatal().Fields(keyvals).Msg(msg)
	}
}

func (zl *zeroLogger) PanicKV(msg string, keyvals map[string]interface{}) {
	for _, l := range zl.loggers {
		l.Panic().Fields(keyvals).Msg(msg)
	}
}

// buildZeroLogger builds a zeroLogger out of the specified "Config".
func buildZeroLogger(config *Config) Logger {

	if !config.Enabled {
		return nonLogger{}
	}

	zerolog.CallerSkipFrameCount = 4

	zerologger := &zeroLogger{make(map[string]*zerolog.Logger, len(config.Appenders))}

	zerolog.TimeFieldFormat = config.TimeFieldFormat
	zerolog.SetGlobalLevel(parseLogLevel(config.GlobalLogLevel))
	zerologger.SetOutputFlags(config.OutputFlags)

	for s, a := range config.Appenders {
		if !a.Enabled {
			continue
		}
		var logger zerolog.Logger
		context := zerolog.New(a.Output).Level(parseLogLevel(a.LogLevel)).With().Timestamp()
		if a.ShowCaller {
			context = context.Caller()
		}
		if a.ShowHostname {
			hostname, err := os.Hostname()
			if err != nil {
				panic(err)
			}
			context = context.Str(globalOutputFlags.HostnameFieldName, hostname)
		}
		logger = context.Logger()
		if a.Format == TextFmt {
			if a.Output == os.Stdout || a.Output == os.Stderr {
				logger = logger.Output(zerolog.ConsoleWriter{Out: a.Output})
			} else {
				logger = logger.Output(zerolog.ConsoleWriter{Out: a.Output, NoColor: true}).With().Timestamp().Logger()
			}
		}
		zerologger.loggers[s] = &logger
	}
	return zerologger
}

// parseLogLevel matches our log-level with zerolog's Level.
func parseLogLevel(level Level) zerolog.Level {
	switch level {
	case DebugLevel:
		return zerolog.DebugLevel
	case InfoLevel:
		return zerolog.InfoLevel
	case WarnLevel:
		return zerolog.WarnLevel
	case ErrorLevel:
		return zerolog.ErrorLevel
	case FatalLevel:
		return zerolog.FatalLevel
	case PanicLevel:
		return zerolog.PanicLevel
	case Disabled:
		return zerolog.Disabled
	}
	return zerolog.Disabled
}
