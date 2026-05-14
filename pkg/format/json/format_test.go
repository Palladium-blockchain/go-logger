package json_test

import (
	"encoding/json"
	"testing"

	"github.com/Palladium-blockchain/go-logger/pkg/core"
	logjson "github.com/Palladium-blockchain/go-logger/pkg/format/json"
)

func TestFormatterFormat(t *testing.T) {
	formatter := &logjson.Formatter{}

	got, err := formatter.Format(core.Record{
		Level:   core.Info,
		Message: "node started",
		Fields: []core.Field{
			{Key: "height", Value: 42},
			{Key: "enabled", Value: true},
			{Key: "role", Value: "validator"},
		},
	})
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := `{"level":"info","message":"node started","fields":[{"key":"height","value":42},{"key":"enabled","value":true},{"key":"role","value":"validator"}]}` + "\n"
	if string(got) != want {
		t.Fatalf("Format() = %q, want %q", got, want)
	}
}

func TestFormatterFormatReturnsMarshalError(t *testing.T) {
	formatter := &logjson.Formatter{}

	_, err := formatter.Format(core.Record{
		Level:   core.Error,
		Message: "not serializable",
		Fields: []core.Field{
			{Key: "value", Value: make(chan int)},
		},
	})
	if err == nil {
		t.Fatal("Format returned nil error for unsupported field value")
	}
	if _, ok := err.(*json.UnsupportedTypeError); !ok {
		t.Fatalf("Format returned error %T, want *json.UnsupportedTypeError", err)
	}
}
