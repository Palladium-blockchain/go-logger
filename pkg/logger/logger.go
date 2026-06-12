package logger

import (
	"io"
	"os"
	"sync/atomic"

	"github.com/Palladium-blockchain/go-logger/pkg/core"
	"github.com/Palladium-blockchain/go-logger/pkg/encoder/plain"
	"github.com/Palladium-blockchain/go-logger/pkg/iox"
)

type Logger interface {
	Debug(msg string, fields ...core.Field)
	Info(msg string, fields ...core.Field)
	Warn(msg string, fields ...core.Field)
	Error(msg string, fields ...core.Field)
	SetLevel(level core.Level)
}

type Option func(logger *logger)

func WithEncoder(encoder core.Encoder) Option {
	return func(logger *logger) {
		logger.encoder = encoder
	}
}

func WithLockingWriter(writer io.Writer) Option {
	return func(logger *logger) {
		logger.writer = iox.NewLockingWriter(writer)
	}
}

func WithWriter(writer core.Writer) Option {
	return func(logger *logger) {
		logger.writer = writer
	}
}

func WithErrorHandler(handlerFn func(error)) Option {
	return func(logger *logger) {
		logger.errorHandler = handlerFn
	}
}

func WithLogLevel(level core.Level) Option {
	return func(logger *logger) {
		logger.SetLevel(level)
	}
}

func New(opts ...Option) Logger {
	l := &logger{
		writer:  iox.NewLockingWriter(os.Stdout),
		encoder: plain.NewEncoder(),
	}
	l.SetLevel(core.Debug)

	for _, opt := range opts {
		opt(l)
	}

	return l
}

type logger struct {
	encoder      core.Encoder
	writer       core.Writer
	logLevel     atomic.Int32
	errorHandler func(error)
}

func (l *logger) Debug(msg string, fields ...core.Field) {
	l.write(core.Record{
		Level:   core.Debug,
		Message: msg,
		Fields:  fields,
	})
}

func (l *logger) Info(msg string, fields ...core.Field) {
	l.write(core.Record{
		Level:   core.Info,
		Message: msg,
		Fields:  fields,
	})
}

func (l *logger) Warn(msg string, fields ...core.Field) {
	l.write(core.Record{
		Level:   core.Warn,
		Message: msg,
		Fields:  fields,
	})
}

func (l *logger) Error(msg string, fields ...core.Field) {
	l.write(core.Record{
		Level:   core.Error,
		Message: msg,
		Fields:  fields,
	})
}

func (l *logger) SetLevel(level core.Level) {
	priority, ok := levelPriority(level)
	if !ok {
		return
	}

	l.logLevel.Store(priority)
}

func (l *logger) write(record core.Record) {
	if !l.shouldWrite(record.Level) {
		return
	}

	if err := l.encoder.Encode(l.writer, record); err != nil && l.errorHandler != nil {
		l.errorHandler(err)
	}
}

func (l *logger) shouldWrite(level core.Level) bool {
	recordPriority, ok := levelPriority(level)
	if !ok {
		return true
	}

	return recordPriority >= l.logLevel.Load()
}

func levelPriority(level core.Level) (int32, bool) {
	switch level {
	case core.Debug:
		return 0, true
	case core.Info:
		return 1, true
	case core.Warn:
		return 2, true
	case core.Error:
		return 3, true
	default:
		return 0, false
	}
}
