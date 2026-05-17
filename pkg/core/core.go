package core

import "io"

type Level string

const (
	Debug Level = "debug"
	Info  Level = "info"
	Warn  Level = "warning"
	Error Level = "error"
)

type Field struct {
	Key   string
	Value any
}

type Record struct {
	Level   Level
	Message string
	Fields  []Field
}

type Encoder interface {
	Encode(w io.Writer, record Record) error
}

type Writer interface {
	Write(buf []byte) (int, error)
}
