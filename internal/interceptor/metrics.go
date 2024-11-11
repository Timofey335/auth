package interceptor

import (
	"context"

	"google.golang.org/grpc"

	"github.com/Timofey335/auth/internal/metric"
)

func MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	metric.IncRequestCounter()

	return handler(ctx, req)
}
