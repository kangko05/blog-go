package logger

import (
	"fmt"
	"log"
)

type Logger interface {
	Info(msg string, data map[string]any)
	Debug(msg string, data map[string]any)
	Warn(msg string, data map[string]any)
	Error(msg string, data map[string]any)
	Write(flag LogLevel, msg string, data map[string]any)
}

type ConsoleLogger struct{}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (l *ConsoleLogger) Info(msg string, data map[string]any) {
	l.Write(LogInfo, msg, data)
}

func (l *ConsoleLogger) Debug(msg string, data map[string]any) {
	l.Write(LogDebug, msg, data)
}

func (l *ConsoleLogger) Warn(msg string, data map[string]any) {
	l.Write(LogWarn, msg, data)
}

func (l *ConsoleLogger) Error(msg string, data map[string]any) {
	l.Write(LogError, msg, data)
}

func (l *ConsoleLogger) Write(flag LogLevel, msg string, data map[string]any) {
	log.Printf("[%s] %s\n", flag, msg)

	for k, v := range data {
		fmt.Printf("\t%s: %v\n", k, v)
	}
}
