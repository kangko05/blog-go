package router

import (
	"blog-go/logger"
	"time"
)

// "requestIP": ctx.ClientIP(),
// "status":    ctx.Writer.Status(),
// "duration":  duration,

type HandlerResult struct {
	LogLevel      logger.LogLevel
	RequestIP     string
	Message       string
	RequestURL    string
	RequestMethod string
	Duration      time.Duration
	Status        int
}

func (hr *HandlerResult) MapResult() map[string]any {
	return map[string]any{
		"logLevel":      hr.LogLevel,
		"requestIP":     hr.RequestIP,
		"requestURL":    hr.RequestURL,
		"requestMethod": hr.RequestMethod,
		"duratin":       hr.Duration,
		"status":        hr.Status,
		"message":       hr.Message,
	}
}

func (hr *HandlerResult) SetHttpStatus(logLevel logger.LogLevel, msg string, statusCode int) {
	hr.LogLevel = logLevel
	hr.Status = statusCode
	hr.Message = msg
}
