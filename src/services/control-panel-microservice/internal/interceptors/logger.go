package interceptors

import (
	"context"

	"control-panel/pkg/api/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func ContextWithLogger(l logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		l.Info(ctx, "request started", zap.String("method", info.FullMethod))
		return handler(ctx, req)
	}
}
