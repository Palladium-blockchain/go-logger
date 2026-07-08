package logger

import (
	"github.com/Palladium-blockchain/go-logger/pkg/core"
)

type noopLogger struct{}

func NewNoopLogger() Logger {
	return &noopLogger{}
}

func (l *noopLogger) Debug(_ string, _ ...core.Field) {}

func (l *noopLogger) Info(_ string, _ ...core.Field) {}

func (l *noopLogger) Warn(_ string, _ ...core.Field) {}

func (l *noopLogger) Error(_ string, _ ...core.Field) {}

func (l *noopLogger) SetLevel(_ß core.Level) {}
