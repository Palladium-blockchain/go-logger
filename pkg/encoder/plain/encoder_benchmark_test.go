package plain_test

import (
	"io"
	"testing"

	"github.com/Palladium-blockchain/go-logger/pkg/core"
	"github.com/Palladium-blockchain/go-logger/pkg/encoder/plain"
)

var plainRecord = core.Record{
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
	encoder := plain.NewEncoder()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := encoder.Encode(io.Discard, plainRecord); err != nil {
			b.Fatal(err)
		}
	}
}
