package logger

import (
	"context"

	"github.com/pkg/errors"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

// InterceptorLogger is a wrapper function for grpc logger middleware.
func InterceptorLogger() logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		switch lvl {
		case logging.LevelDebug:
			Debug(msg, fields...)
		case logging.LevelInfo:
			Info(msg, fields...)
		case logging.LevelWarn:
			Warn(msg, fields...)
		case logging.LevelError:
			Error(msg, errors.New("grpc error"), fields...)
		}
	})
}
