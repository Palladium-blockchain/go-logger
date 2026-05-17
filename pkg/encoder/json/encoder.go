package json

import (
	"encoding/json"
	"io"

	"github.com/Palladium-blockchain/go-logger/pkg/core"
)

type Field struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type Encoder struct{}

func (e *Encoder) Encode(w io.Writer, record core.Record) error {
	// TODO: rewrite to manual json encoding
	entry := make(map[string]any, len(record.Fields)+2)

	entry["level"] = string(record.Level)
	entry["message"] = record.Message

	for _, field := range record.Fields {
		switch field.Key {
		case "level", "message":
			entry["field."+field.Key] = field.Value
		default:
			entry[field.Key] = field.Value
		}
	}

	return json.NewEncoder(w).Encode(entry)
}
