package log

import (
	"context"
	"log/slog"
)

const (
	defaultErrorField = "error"
	defaultDataField  = "data"
)

type slogLogger struct {
	logger *slog.Logger
}

func NewLogger(opts ...Option) Logger {
	options := newDefaultOptions()

	for _, opt := range opts {
		opt(options)
	}

	var handler slog.Handler
	if options.handler != nil {
		handler = options.handler
	} else {
		handler = slog.NewJSONHandler(options.output, &slog.HandlerOptions{ //nolint:exhaustruct
			Level: options.level,
		})
	}

	return &slogLogger{
		logger: slog.New(handler),
	}
}

func NewNopLogger() Logger {
	return &slogLogger{
		logger: slog.New(slog.DiscardHandler),
	}
}

func (l *slogLogger) Error(ctx context.Context, msg string, kv ...KeyValue) {
	l.logger.ErrorContext(ctx, msg, l.toArgs(kv)...)
}

func (l *slogLogger) Warn(ctx context.Context, msg string, kv ...KeyValue) {
	l.logger.WarnContext(ctx, msg, l.toArgs(kv)...)
}

func (l *slogLogger) Info(ctx context.Context, msg string, kv ...KeyValue) {
	l.logger.InfoContext(ctx, msg, l.toArgs(kv)...)
}

func (l *slogLogger) Debug(ctx context.Context, msg string, kv ...KeyValue) {
	l.logger.DebugContext(ctx, msg, l.toArgs(kv)...)
}

func (l *slogLogger) With(kv ...KeyValue) Logger {
	return &slogLogger{
		logger: l.logger.With(l.toArgs(kv)...),
	}
}

func (l *slogLogger) WithErr(err error) Logger {
	return &slogLogger{
		logger: l.logger.With(slog.String(defaultErrorField, err.Error())),
	}
}

func (l *slogLogger) toArgs(kvs []KeyValue) []any {
	if len(kvs) == 0 {
		return nil
	}

	args := make([]any, 0, len(kvs)*2)
	for _, kv := range kvs {
		args = append(args, kv.Key, kv.Value)
	}

	return []any{slog.Group(defaultDataField, args...)}
}
