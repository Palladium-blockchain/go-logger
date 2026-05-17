package logger

import (
	"io"
	"os"

	"github.com/Palladium-blockchain/go-logger/pkg/core"
	"github.com/Palladium-blockchain/go-logger/pkg/encoder/plain"
	"github.com/Palladium-blockchain/go-logger/pkg/iox"
)

type Logger interface {
	Debug(msg string, fields ...core.Field)
	Info(msg string, fields ...core.Field)
	Warn(msg string, fields ...core.Field)
	Error(msg string, fields ...core.Field)
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

func New(opts ...Option) Logger {
	l := &logger{
		writer:  iox.NewLockingWriter(os.Stdout),
		encoder: plain.NewEncoder(),
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

type logger struct {
	encoder      core.Encoder
	writer       core.Writer
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

func (l *logger) write(record core.Record) {
	if err := l.encoder.Encode(l.writer, record); err != nil && l.errorHandler != nil {
		l.errorHandler(err)
	}
}
