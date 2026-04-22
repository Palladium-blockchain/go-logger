package logger

import (
	"fmt"
	"log"
)

type Logger interface {
	Debug(msg string, fields ...*Field)
	Info(msg string, fields ...*Field)
	Warn(msg string, fields ...*Field)
	Error(msg string, fields ...*Field)
}

func New() Logger {
	return &logger{}
}

type logger struct {
}

func (_ *logger) Debug(msg string, fields ...*Field) {
	Debug(msg, fields...)
}

func (_ *logger) Info(msg string, fields ...*Field) {
	Info(msg, fields...)
}

func (_ *logger) Warn(msg string, fields ...*Field) {
	Warn(msg, fields...)
}

func (_ *logger) Error(msg string, fields ...*Field) {
	Error(msg, fields...)
}

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

type Field struct {
	Key   string
	Value any
}

func WithField(key string, value any) *Field {
	return &Field{Key: key, Value: value}
}

func WithError(err error) *Field {
	return &Field{Key: "error", Value: err}
}

func Debug(msg string, fields ...*Field) {
	l(DEBUG, msg, fields...)
}

func Info(msg string, fields ...*Field) {
	l(INFO, msg, fields...)
}

func Warn(msg string, fields ...*Field) {
	l(WARNING, msg, fields...)
}

func Error(msg string, fields ...*Field) {
	l(ERROR, msg, fields...)
}

func l(levelLevel LogLevel, msg string, fields ...*Field) {
	rec := msg
	for _, field := range fields {
		rec += " " + field.Key + ": " + fmt.Sprintf("%v", field.Value)
	}
	log.Println(rec)
}
