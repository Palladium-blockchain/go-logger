package json_test

import (
	"io"
	"testing"

	"github.com/Palladium-blockchain/go-logger/pkg/core"
	logjson "github.com/Palladium-blockchain/go-logger/pkg/encoder/json"
)

var jsonRecord = core.Record{
	Level:   core.Info,
	Message: "node started",
	Fields: []core.Field{
		{Key: "height", Value: 42},
		{Key: "role", Value: "validator"},
		{Key: "healthy", Value: true},
		{Key: "latency", Value: 12.35},
	},
}

func BenchmarkEncoder(b *testing.B) {
	encoder := &logjson.Encoder{}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := encoder.Encode(io.Discard, jsonRecord); err != nil {
			b.Fatal(err)
		}
	}
}
