package logger

import (
	"context"

	l "control-panel/pkg/logger"
	"go.uber.org/zap"
)

const (
	LoggerKey = "logger"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Sync() error
}

func New(serviceName string) Logger {
	logger, err := l.NewLogger(serviceName)
	if err != nil {
		panic(err)
	}
	return logger
}

func GetLoggerFromCtx(ctx context.Context) Logger {
	return ctx.Value(LoggerKey).(Logger)
}
