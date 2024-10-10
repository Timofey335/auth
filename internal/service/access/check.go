package access

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	authPrefix = "Bearer "
)

func (*serv) Check(ctx context.Context, s string) (*emptypb.Empty, error) {

	return &emptypb.Empty{}, nil
}

func authHeader(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return "", errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return "", errors.New("invalid authorization header format")
	}

	return authHeader[0], nil
}
