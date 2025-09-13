package logger

type LogLevel string

const (
	LogInfo  LogLevel = "info"
	LogDebug LogLevel = "debug"
	LogWarn  LogLevel = "warn"
	LogError LogLevel = "error"
)
