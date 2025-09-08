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
}

type ConsoleLogger struct{}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (l *ConsoleLogger) Info(msg string, data map[string]any) {
	l.write("INFO", msg, data)
}

func (l *ConsoleLogger) Debug(msg string, data map[string]any) {
	l.write("DEBUG", msg, data)
}

func (l *ConsoleLogger) Warn(msg string, data map[string]any) {
	l.write("WARN", msg, data)
}

func (l *ConsoleLogger) Error(msg string, data map[string]any) {
	l.write("ERROR", msg, data)
}

func (l *ConsoleLogger) write(flag, msg string, data map[string]any) {
	log.Printf("[%s] %s\n", flag, msg)

	for k, v := range data {
		fmt.Printf("\t%s: %v\n", k, v)
	}
}
