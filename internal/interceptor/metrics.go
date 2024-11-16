package interceptor

import (
	"context"

	"github.com/Timofey335/auth/internal/metric"
	"google.golang.org/grpc"
)

func MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	metric.IncRequestCounter()

	res, err := handler(ctx, req)

	if err != nil {
		metric.IncResponseCounter("error", info.FullMethod)
	} else {
		metric.IncResponseCounter("success", info.FullMethod)
	}

	return res, err
}
