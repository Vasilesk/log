package log

import (
	"io"
	"log/slog"
	"os"
)

type Option func(*options)

type options struct {
	output  io.Writer
	level   slog.Level
	handler slog.Handler
}

func WithOutput(w io.Writer) Option {
	return func(o *options) {
		o.output = w
	}
}

func WithLevel(l slog.Level) Option {
	return func(o *options) {
		o.level = l
	}
}

func WithHandler(h slog.Handler) Option {
	return func(o *options) {
		o.handler = h
	}
}

func newDefaultOptions() *options {
	return &options{
		output:  os.Stdout,
		level:   slog.LevelInfo,
		handler: nil,
	}
}
