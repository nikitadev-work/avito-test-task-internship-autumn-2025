package logger

import (
	"log/slog"
	"os"
	"strings"
)

type LoggerInterface interface {
	Debug(msg string, params map[string]any)
	Info(msg string, params map[string]any)
	Warn(msg string, params map[string]any)
	Error(msg string, params map[string]any)
}

type Logger struct {
	logger *slog.Logger
}

var _ LoggerInterface = (*Logger)(nil)

func NewLogger(levelStr string, params map[string]any) *Logger {
	var level slog.Level

	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		level = slog.LevelDebug
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	})
	l := slog.New(handler)
	l = l.With(paramsToAny(params)...)

	return &Logger{
		logger: l,
	}
}

func (l *Logger) Debug(msg string, params map[string]any) {
	l.logger.Debug(msg, paramsToAny(params)...)
}

func (l *Logger) Info(msg string, params map[string]any) {
	l.logger.Info(msg, paramsToAny(params)...)
}

func (l *Logger) Warn(msg string, params map[string]any) {
	l.logger.Warn(msg, paramsToAny(params)...)
}

func (l *Logger) Error(msg string, params map[string]any) {
	l.logger.Error(msg, paramsToAny(params)...)
}

func (l *Logger) With(params map[string]any) *Logger {
	return &Logger{
		logger: l.logger.With(paramsToAny(params)...),
	}
}

func paramsToAny(params map[string]any) []any {
	if len(params) == 0 {
		return nil
	}

	res := make([]any, 0, len(params)*2)

	for k, v := range params {
		res = append(res, k, v)
	}

	return res
}
