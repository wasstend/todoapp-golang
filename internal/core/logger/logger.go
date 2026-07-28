package core_logger

import (
	"context"
)

type Logger interface {
	Debug(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	With(fields ...Field) Logger
	Close()
}

type loggerContextKey struct{}

var (
	key = loggerContextKey{}
)

func IntoContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx,
		key,
		logger,
	)
}

func FromContext(ctx context.Context) Logger {
	log, ok := ctx.Value(key).(Logger)
	if !ok {
		panic("no logger in context")
	}
	return log
}
