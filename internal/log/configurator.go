package log

import (
	"log/slog"
	"os"
)

type logParams struct {
	level  slog.Level
	format string
}

func WithLevel(level string) func(p *logParams) {
	return func(p *logParams) {
		p.level = parseLevel(level)
	}
}

func WithFormat(format string) func(p *logParams) {
	return func(p *logParams) {
		p.format = format
	}
}

func defaultParams() logParams {
	return logParams{
		level:  slog.LevelInfo,
		format: "json",
	}
}

func newHandlerOptions(level slog.Level) *slog.HandlerOptions {
	return &slog.HandlerOptions{
		Level:       level,
		AddSource:   false,
		ReplaceAttr: nil,
	}
}

func newHandler(format string, options *slog.HandlerOptions) slog.Handler {
	if format == "text" {
		return slog.NewTextHandler(os.Stdout, options)
	} else {
		return slog.NewJSONHandler(os.Stdout, options)
	}
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
