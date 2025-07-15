package middleware

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime/debug"
)

func NewGrpcPanicRecoveryHandler(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(func(p any) (err error) {
		logger.Error("recovered from panic", "panic", p, "stack", debug.Stack())
		return status.Errorf(codes.Internal, "%s", p)
	}))
}
