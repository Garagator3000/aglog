package log

import (
	"log/slog"
)

func Error(err error) slog.Attr {
	return slog.String("error", err.Error())
}

func String(key string, value string) slog.Attr {
	return slog.String(key, value)
}

func Int(key string, value int) slog.Attr {
	return slog.Int(key, value)
}

func Int64(key string, value int64) slog.Attr {
	return slog.Int64(key, value)
}

func Bool(key string, value bool) slog.Attr {
	return slog.Bool(key, value)
}
