package core

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

type Formatter interface {
	Format(record Record) ([]byte, error)
}

type Writer interface {
	Write(buf []byte) (int, error)
}
