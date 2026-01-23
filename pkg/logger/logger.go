package logger

import (
	"fmt"
	"log"
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

func Info(msg string, fields ...*Field) {
	l(msg, fields...)
}

func Error(msg string, fields ...*Field) {
	l(msg, fields...)
}

func l(msg string, fields ...*Field) {
	rec := msg
	for _, field := range fields {
		rec += " " + field.Key + ": " + fmt.Sprintf("%v", field.Value)
	}
	log.Println(rec)
}
