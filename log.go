package log

import (
	"context"

	"github.com/samber/lo"
)

type Logger interface {
	Error(ctx context.Context, msg string, kv ...KeyValue)
	Warn(ctx context.Context, msg string, kv ...KeyValue)
	Info(ctx context.Context, msg string, kv ...KeyValue)
	Debug(ctx context.Context, msg string, kv ...KeyValue)
	With(kv ...KeyValue) Logger
	WithErr(err error) Logger
}

type KeyValue = lo.Entry[string, any]

func KV(key string, value any) KeyValue {
	return KeyValue{
		Key:   key,
		Value: value,
	}
}
