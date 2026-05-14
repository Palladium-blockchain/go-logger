package json

import (
	"encoding/json"

	"github.com/Palladium-blockchain/go-logger/pkg/core"
)

type Record struct {
	Level   string  `json:"level"`
	Message string  `json:"message"`
	Fields  []Field `json:"fields"`
}

type Field struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type Formatter struct{}

func (f *Formatter) Format(record core.Record) ([]byte, error) {
	fields := make([]Field, len(record.Fields))
	for i, field := range record.Fields {
		fields[i] = Field{
			Key:   field.Key,
			Value: field.Value,
		}
	}

	jsonRecord := Record{
		Level:   string(record.Level),
		Message: record.Message,
		Fields:  fields,
	}

	buf, err := json.Marshal(jsonRecord)
	if err != nil {
		return nil, err
	}

	buf = append(buf, '\n')
	return buf, nil
}
