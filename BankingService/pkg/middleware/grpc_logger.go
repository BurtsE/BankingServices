package middleware

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/sirupsen/logrus"
)

func NewInterceptorLogger(l *logrus.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(l.Level, fields...)
	})
}
