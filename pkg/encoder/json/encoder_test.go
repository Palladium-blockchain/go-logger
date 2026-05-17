package json_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/Palladium-blockchain/go-logger/pkg/core"
	logjson "github.com/Palladium-blockchain/go-logger/pkg/encoder/json"
)

func TestEncoderEncode(t *testing.T) {
	encoder := &logjson.Encoder{}
	var buf bytes.Buffer

	err := encoder.Encode(&buf, core.Record{
		Level:   core.Info,
		Message: "node started",
		Fields: []core.Field{
			{Key: "height", Value: 42},
			{Key: "enabled", Value: true},
			{Key: "role", Value: "validator"},
		},
	})
	if err != nil {
		t.Fatalf("Encode returned error: %v", err)
	}

	got := decodeObject(t, buf.Bytes())
	want := map[string]any{
		"level":   "info",
		"message": "node started",
		"height":  float64(42),
		"enabled": true,
		"role":    "validator",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Encode() = %#v, want %#v", got, want)
	}
}

func TestEncoderEncodePrefixesReservedFieldKeys(t *testing.T) {
	encoder := &logjson.Encoder{}
	var buf bytes.Buffer

	err := encoder.Encode(&buf, core.Record{
		Level:   core.Info,
		Message: "node started",
		Fields: []core.Field{
			{Key: "level", Value: "user level"},
			{Key: "message", Value: "user message"},
		},
	})
	if err != nil {
		t.Fatalf("Encode returned error: %v", err)
	}

	got := decodeObject(t, buf.Bytes())
	want := map[string]any{
		"level":         "info",
		"message":       "node started",
		"field.level":   "user level",
		"field.message": "user message",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Encode() = %#v, want %#v", got, want)
	}
}

func TestEncoderEncodeReturnsMarshalError(t *testing.T) {
	encoder := &logjson.Encoder{}
	var buf bytes.Buffer

	err := encoder.Encode(&buf, core.Record{
		Level:   core.Error,
		Message: "not serializable",
		Fields: []core.Field{
			{Key: "value", Value: make(chan int)},
		},
	})
	if err == nil {
		t.Fatal("Encode returned nil error for unsupported field value")
	}
	if _, ok := err.(*json.UnsupportedTypeError); !ok {
		t.Fatalf("Encode returned error %T, want *json.UnsupportedTypeError", err)
	}
}

func decodeObject(t *testing.T, b []byte) map[string]any {
	t.Helper()

	var got map[string]any
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("failed to decode output %q: %v", b, err)
	}

	return got
}
