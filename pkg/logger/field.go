package logger

import "github.com/Palladium-blockchain/go-logger/pkg/core"

func WithField(key string, value any) core.Field {
	return core.Field{Key: key, Value: value}
}

func WithError(err error) core.Field {
	return core.Field{Key: "error", Value: err}
}
