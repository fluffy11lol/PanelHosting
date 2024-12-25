package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const loggerKey = "logger"
const RequestID = "requestID"

type Logger struct {
	*zap.Logger
	serviceName string
}

func NewLogger(serviceName string) (*Logger, error) {
	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "timestamp"
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalLevelEncoder

	encoder := zapcore.NewJSONEncoder(config)

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)

	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return &Logger{
		Logger:      logger,
		serviceName: serviceName,
	}, nil
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("service", l.serviceName))
	if reqID := ctx.Value(RequestID); reqID != nil {
		fields = append(fields, zap.String(RequestID, reqID.(string)))
	}
	l.Logger.Info(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("service", l.serviceName))
	if reqID := ctx.Value(RequestID); reqID != nil {
		fields = append(fields, zap.String(RequestID, reqID.(string)))
	}
	l.Logger.Error(msg, fields...)
}

func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("service", l.serviceName))
	if reqID := ctx.Value(RequestID); reqID != nil {
		fields = append(fields, zap.String(RequestID, reqID.(string)))
	}
	l.Logger.Debug(msg, fields...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("service", l.serviceName))
	if reqID := ctx.Value(RequestID); reqID != nil {
		fields = append(fields, zap.String(RequestID, reqID.(string)))
	}
	l.Logger.Warn(msg, fields...)
}

func (l *Logger) WithFields(fields ...zap.Field) *Logger {
	return &Logger{
		Logger:      l.Logger.With(fields...),
		serviceName: l.serviceName,
	}
}

func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

func GetLoggerFromCtx(ctx context.Context) Logger {
	return ctx.Value(loggerKey).(Logger)
}
