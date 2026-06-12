package logger_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/Palladium-blockchain/go-logger/pkg/core"
	"github.com/Palladium-blockchain/go-logger/pkg/logger"
)

type recordingEncoder struct {
	records []core.Record
	output  []byte
}

func (e *recordingEncoder) Encode(w io.Writer, record core.Record) error {
	e.records = append(e.records, record)
	_, err := w.Write(e.output)
	return err
}

func TestLoggerMethodsEncodeExpectedRecords(t *testing.T) {
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
			encoder := &recordingEncoder{output: []byte("formatted log")}
			var writer bytes.Buffer
			log := logger.New(
				logger.WithEncoder(encoder),
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
			if got, want := len(encoder.records), 1; got != want {
				t.Fatalf("encoder got %d records, want %d", got, want)
			}

			want := core.Record{
				Level:   tt.level,
				Message: "node started",
				Fields:  fields,
			}
			if !reflect.DeepEqual(encoder.records[0], want) {
				t.Fatalf("encoder record = %#v, want %#v", encoder.records[0], want)
			}
		})
	}
}

func TestNewUsesConfiguredWriter(t *testing.T) {
	encoder := &recordingEncoder{output: []byte("custom output")}
	var writer bytes.Buffer
	log := logger.New(
		logger.WithEncoder(encoder),
		logger.WithLockingWriter(&writer),
	)

	log.Info("configured writer")

	if got, want := writer.String(), "custom output"; got != want {
		t.Fatalf("writer got %q, want %q", got, want)
	}
}

func TestWithLogLevelFiltersLowerPriorityRecords(t *testing.T) {
	encoder := &recordingEncoder{output: []byte("record\n")}
	var writer bytes.Buffer
	log := logger.New(
		logger.WithEncoder(encoder),
		logger.WithWriter(&writer),
		logger.WithLogLevel(core.Warn),
	)

	log.Debug("debug")
	log.Info("info")
	log.Warn("warn")
	log.Error("error")

	if got, want := writer.String(), "record\nrecord\n"; got != want {
		t.Fatalf("writer got %q, want %q", got, want)
	}
	if got, want := len(encoder.records), 2; got != want {
		t.Fatalf("encoder got %d records, want %d", got, want)
	}

	wantLevels := []core.Level{core.Warn, core.Error}
	for i, want := range wantLevels {
		if got := encoder.records[i].Level; got != want {
			t.Fatalf("record %d level = %q, want %q", i, got, want)
		}
	}
}

func TestWithLogLevelIgnoresUnknownLevel(t *testing.T) {
	encoder := &recordingEncoder{output: []byte("record")}
	var writer bytes.Buffer
	log := logger.New(
		logger.WithEncoder(encoder),
		logger.WithWriter(&writer),
		logger.WithLogLevel(core.Level("unknown")),
	)

	log.Debug("debug")

	if got, want := len(encoder.records), 1; got != want {
		t.Fatalf("encoder got %d records, want %d", got, want)
	}
}
