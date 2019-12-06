package log

// Interface assertions
var _ Logger = (*nonLogger)(nil)

// nonLogger is used when we disable logging.
type nonLogger struct{}

func (nonLogger) SetTimeFieldFormat(format string) {}

func (nonLogger) SetOutputFlags(flags *OutputFlags) {}

func (nonLogger) SetGlobalLogLevel(level Level) {}

func (nonLogger) Debug(msg string) {}

func (nonLogger) Info(msg string) {}

func (nonLogger) Warn(msg string) {}

func (nonLogger) Error(msg string) {}

func (nonLogger) Fatal(msg string) {}

func (nonLogger) Panic(msg string) {}

func (nonLogger) DebugKV(msg string, keyvals map[string]interface{}) {}

func (nonLogger) InfoKV(msg string, keyvals map[string]interface{}) {}

func (nonLogger) WarnKV(msg string, keyvals map[string]interface{}) {}

func (nonLogger) ErrorKV(msg string, keyvals map[string]interface{}) {}

func (nonLogger) FatalKV(msg string, keyvals map[string]interface{}) {}

func (nonLogger) PanicKV(msg string, keyvals map[string]interface{}) {}
