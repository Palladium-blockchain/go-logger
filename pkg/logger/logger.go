package logger

import (
	"io"
	"os"

	"github.com/Palladium-blockchain/go-logger/pkg/core"
	"github.com/Palladium-blockchain/go-logger/pkg/format/plain"
	"github.com/Palladium-blockchain/go-logger/pkg/iox"
)

type Logger interface {
	Debug(msg string, fields ...core.Field)
	Info(msg string, fields ...core.Field)
	Warn(msg string, fields ...core.Field)
	Error(msg string, fields ...core.Field)
}

type Option func(logger *logger)

func WithFormatter(formatter core.Formatter) Option {
	return func(logger *logger) {
		logger.formatter = formatter
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

func New(opts ...Option) Logger {
	l := &logger{
		writer:    iox.NewLockingWriter(os.Stdout),
		formatter: &plain.Formatter{},
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

type logger struct {
	formatter core.Formatter
	writer    core.Writer
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
	// TODO: implement error handling
	bytes, _ := l.formatter.Format(record)
	_, _ = l.writer.Write(bytes)
}
