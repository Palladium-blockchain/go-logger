package logger_test

import (
	"io"
	"testing"

	logjson "github.com/Palladium-blockchain/go-logger/pkg/encoder/json"
	"github.com/Palladium-blockchain/go-logger/pkg/encoder/plain"
	"github.com/Palladium-blockchain/go-logger/pkg/logger"
)

var benchmarkFields = []struct {
	name string
	log  logger.Logger
}{
	{
		name: "PlainDiscardWriter",
		log: logger.New(
			logger.WithEncoder(plain.NewEncoder()),
			logger.WithWriter(io.Discard),
		),
	},
	{
		name: "PlainLockingWriter",
		log: logger.New(
			logger.WithEncoder(plain.NewEncoder()),
			logger.WithLockingWriter(io.Discard),
		),
	},
	{
		name: "JSONDiscardWriter",
		log: logger.New(
			logger.WithEncoder(&logjson.Encoder{}),
			logger.WithWriter(io.Discard),
		),
	},
}

var fields = []loggerField{
	{key: "height", value: 42},
	{key: "role", value: "validator"},
	{key: "healthy", value: true},
	{key: "latency", value: 12.35},
}

type loggerField struct {
	key   string
	value any
}

func BenchmarkLoggerInfo(b *testing.B) {
	for _, bench := range benchmarkFields {
		b.Run(bench.name+"/NoFields", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				bench.log.Info("node started")
			}
		})

		b.Run(bench.name+"/WithFields", func(b *testing.B) {
			height := logger.WithField(fields[0].key, fields[0].value)
			role := logger.WithField(fields[1].key, fields[1].value)
			healthy := logger.WithField(fields[2].key, fields[2].value)
			latency := logger.WithField(fields[3].key, fields[3].value)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				bench.log.Info("node started", height, role, healthy, latency)
			}
		})
	}
}

func BenchmarkLoggerInfoParallel(b *testing.B) {
	for _, bench := range benchmarkFields {
		b.Run(bench.name, func(b *testing.B) {
			height := logger.WithField("height", 42)
			role := logger.WithField("role", "validator")
			healthy := logger.WithField("healthy", true)
			latency := logger.WithField("latency", 12.35)

			b.ReportAllocs()
			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					bench.log.Info("node started", height, role, healthy, latency)
				}
			})
		})
	}
}
