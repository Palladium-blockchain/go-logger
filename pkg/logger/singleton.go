package logger

import (
	"github.com/Palladium-blockchain/go-logger/pkg/core"
)

var stdLogger Logger = New()

func Debug(msg string, fields ...core.Field) {
	stdLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...core.Field) {
	stdLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...core.Field) {
	stdLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...core.Field) {
	stdLogger.Error(msg, fields...)
}
