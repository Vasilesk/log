package log

import "github.com/samber/lo"

type Logger interface {
	Error(msg string, kv ...KV)
	Warn(msg string, kv ...KV)
	Info(msg string, kv ...KV)
	Debug(msg string, kv ...KV)
	With(kv ...KV) Logger
	WithErr(err error) Logger
}

type KV = lo.Entry[string, any]
