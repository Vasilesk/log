package log

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSlogLogger(t *testing.T) {
	tests := []struct {
		name          string
		levelFunc     slog.Level
		logFunc       func(Logger, context.Context, string, ...KeyValue)
		msg           string
		kv            []KeyValue
		expectedLevel string
		expectedData  map[string]any
	}{
		{
			name:          "info level with no kv",
			levelFunc:     slog.LevelInfo,
			logFunc:       Logger.Info,
			msg:           "test info message",
			kv:            nil,
			expectedLevel: "INFO",
		},
		{
			name:      "error level with kv",
			levelFunc: slog.LevelError,
			logFunc:   Logger.Error,
			msg:       "test error message",
			kv: []KeyValue{
				{Key: "key1", Value: "value1"},
				{Key: "key2", Value: 42},
			},
			expectedLevel: "ERROR",
			expectedData: map[string]any{
				"key1": "value1",
				"key2": float64(42), // JSON numbers are float64
			},
		},
		{
			name:      "debug level with single kv",
			levelFunc: slog.LevelDebug,
			logFunc:   Logger.Debug,
			msg:       "test debug message",
			kv: []KeyValue{
				{Key: "debug_key", Value: true},
			},
			expectedLevel: "DEBUG",
			expectedData: map[string]any{
				"debug_key": true,
			},
		},
		{
			name:      "warn level with kv",
			levelFunc: slog.LevelWarn,
			logFunc:   Logger.Warn,
			msg:       "test warn message",
			kv: []KeyValue{
				{Key: "warn_code", Value: "WARN001"},
			},
			expectedLevel: "WARN",
			expectedData: map[string]any{
				"warn_code": "WARN001",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLogger(
				WithOutput(&buf),
				WithLevel(tt.levelFunc),
			)

			tt.logFunc(logger, context.Background(), tt.msg, tt.kv...)

			var actual map[string]any
			err := json.Unmarshal(buf.Bytes(), &actual)
			require.NoError(t, err)

			assert.Equal(t, tt.msg, actual["msg"])
			assert.Equal(t, tt.expectedLevel, actual["level"])

			if tt.expectedData != nil {
				data, ok := actual["data"].(map[string]any)
				require.True(t, ok)
				assert.Equal(t, tt.expectedData, data)
			}
		})
	}
}

func TestSlogLoggerWith(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(WithOutput(&buf))

	withLogger := logger.With(KeyValue{Key: "request_id", Value: "123"})
	withLogger.Info(context.Background(), "test message")

	var actual map[string]any
	err := json.Unmarshal(buf.Bytes(), &actual)
	require.NoError(t, err)

	data, ok := actual["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "123", data["request_id"])
}

func TestSlogLoggerWithErr(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(WithOutput(&buf))

	testErr := assert.AnError
	withErrLogger := logger.WithErr(testErr)
	withErrLogger.Info(context.Background(), "test message")

	var actual map[string]any
	err := json.Unmarshal(buf.Bytes(), &actual)
	require.NoError(t, err)

	assert.Equal(t, testErr.Error(), actual["error"])
}

func TestNewNopLogger(t *testing.T) {
	logger := NewNopLogger()

	// Should not panic
	logger.Info(context.Background(), "test message")
	logger.Error(context.Background(), "test message")
	logger.Debug(context.Background(), "test message")
	logger.Warn(context.Background(), "test message")

	withLogger := logger.With(KeyValue{Key: "key", Value: "value"})
	withLogger.Info(context.Background(), "test message")

	withErrLogger := logger.WithErr(assert.AnError)
	withErrLogger.Info(context.Background(), "test message")
}
