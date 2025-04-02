package logger

import (
	"context"
	"path/filepath"
	"strconv"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
)

const (
	Key       = "logger"
	RequestID = "request_id"
)

type Logger struct {
	l *zap.Logger
}

func New(ctx context.Context) (context.Context, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, Key, &Logger{logger})

	return ctx, nil
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	return ctx.Value(Key).(*Logger)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}

	l.l.Info(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}

	l.l.Error(msg, fields...)
}

func InterceptorLogger(l *Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		filds := make([]zap.Field, len(fields))
		for i, f := range fields {
			filds[i] = zap.Any(filepath.Base(strconv.Itoa(i)), f)
		}
		l.l.Info(msg, filds...)
	})
}
