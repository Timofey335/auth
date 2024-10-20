package access

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Timofey335/auth/internal/utils"
)

const (
	authPrefix = "Bearer "
)

// Check - проверяет токен на подлинность и на доступ к эндпоинтам
func (s *serv) Check(ctx context.Context, endpointAddress string) (*emptypb.Empty, error) {
	authHeader, err := authHeader(ctx)
	if err != nil {
		return nil, err
	}

	accessToken := strings.TrimPrefix(authHeader, authPrefix)

	accessTokenSecretKey := s.authConfig.AccessTokenSecretKey()

	claims, err := utils.VerifyToken(accessToken, []byte(accessTokenSecretKey))
	if err != nil {
		return nil, errors.New("access token is invalid")
	}

	accessibleRole, err := s.accessibleRoles(ctx, endpointAddress)
	if err != nil {
		return nil, errors.New("failed to get accessible role")
	}

	if accessibleRole == 0 || accessibleRole == 2 {
		return &emptypb.Empty{}, nil
	}

	if accessibleRole == 1 {
		if claims.Role == 1 {
			return &emptypb.Empty{}, nil
		}
	}

	return nil, errors.New("access denied")
}

func (s *serv) accessibleRoles(ctx context.Context, endpoint string) (int64, error) {
	role, err := s.accessRepository.GetRole(ctx, endpoint)
	if err != nil {
		return 0, err
	}

	return role, nil
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
