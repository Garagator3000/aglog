package log

import (
	"log/slog"
)

type Logger interface {
	Debug(string, ...any)
	Info(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
}

type logParamsCallback func(params *logParams)

type Log struct {
	logger *slog.Logger
}

func NewLog(opts ...logParamsCallback) *Log {
	params := defaultParams()

	for _, opt := range opts {
		opt(&params)
	}

	handlerOptions := newHandlerOptions(params.level)
	handler := newHandler(params.format, handlerOptions)

	return &Log{
		logger: slog.New(handler),
	}
}

func (l *Log) Debug(msg string, attrs ...any) {
	l.logger.Debug(msg, attrs...)
}

func (l *Log) Info(msg string, attrs ...any) {
	l.logger.Info(msg, attrs...)
}

func (l *Log) Warn(msg string, attrs ...any) {
	l.logger.Warn(msg, attrs...)
}

func (l *Log) Error(msg string, attrs ...any) {
	l.logger.Error(msg, attrs...)
}
