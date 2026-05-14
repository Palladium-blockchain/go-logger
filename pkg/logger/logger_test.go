package logger_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/Palladium-blockchain/go-logger/pkg/core"
	"github.com/Palladium-blockchain/go-logger/pkg/logger"
)

type recordingFormatter struct {
	records []core.Record
	output  []byte
}

func (f *recordingFormatter) Format(record core.Record) ([]byte, error) {
	f.records = append(f.records, record)
	return f.output, nil
}

func TestLoggerMethodsFormatExpectedRecords(t *testing.T) {
	tests := []struct {
		name  string
		level core.Level
		log   func(logger.Logger, string, ...core.Field)
	}{
		{
			name:  "Debug",
			level: core.Debug,
			log:   logger.Logger.Debug,
		},
		{
			name:  "Info",
			level: core.Info,
			log:   logger.Logger.Info,
		},
		{
			name:  "Warn",
			level: core.Warn,
			log:   logger.Logger.Warn,
		},
		{
			name:  "Error",
			level: core.Error,
			log:   logger.Logger.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := &recordingFormatter{output: []byte("formatted log")}
			var writer bytes.Buffer
			log := logger.New(
				logger.WithFormatter(formatter),
				logger.WithWriter(&writer),
			)
			fields := []core.Field{
				logger.WithField("height", 42),
				logger.WithField("role", "validator"),
			}

			tt.log(log, "node started", fields...)

			if got, want := writer.String(), "formatted log"; got != want {
				t.Fatalf("writer got %q, want %q", got, want)
			}
			if got, want := len(formatter.records), 1; got != want {
				t.Fatalf("formatter got %d records, want %d", got, want)
			}

			want := core.Record{
				Level:   tt.level,
				Message: "node started",
				Fields:  fields,
			}
			if !reflect.DeepEqual(formatter.records[0], want) {
				t.Fatalf("formatter record = %#v, want %#v", formatter.records[0], want)
			}
		})
	}
}

func TestNewUsesConfiguredWriter(t *testing.T) {
	formatter := &recordingFormatter{output: []byte("custom output")}
	var writer bytes.Buffer
	log := logger.New(
		logger.WithFormatter(formatter),
		logger.WithLockingWriter(&writer),
	)

	log.Info("configured writer")

	if got, want := writer.String(), "custom output"; got != want {
		t.Fatalf("writer got %q, want %q", got, want)
	}
}
