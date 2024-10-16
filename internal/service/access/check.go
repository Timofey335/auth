package access

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Timofey335/auth/internal/model"
	"github.com/Timofey335/auth/internal/utils"
)

const (
	authPrefix = "Bearer "

	refreshTokenSecretKey = "W4/X+LLjehdxptt4YgGFCvMpq5ewptpZZYRHY6A72g0="
	accessTokenSecretKey  = "VqvguGiffXILza1f44TWXowDT4zwf03dtXmqWW4SYyE="
)

var accessibleRoles map[string]int64

// Check - проверяет токен на подлинность и на доступ к эндпоинтам
func (s *serv) Check(ctx context.Context, endpointAddress string) (*emptypb.Empty, error) {
	authHeader, err := authHeader(ctx)
	if err != nil {
		return nil, err
	}

	accessToken := strings.TrimPrefix(authHeader, authPrefix)

	claims, err := utils.VerifyToken(accessToken, []byte(accessTokenSecretKey))
	if err != nil {
		return nil, errors.New("access token is invalid")
	}

	accessibleMap, err := s.accessibleRoles(ctx)
	if err != nil {
		return nil, errors.New("failed to get accessible roles")
	}

	role, ok := accessibleMap[endpointAddress]
	if !ok {
		return &emptypb.Empty{}, nil
	}

	if role == claims.Role {
		return &emptypb.Empty{}, nil
	}

	return nil, errors.New("access denied")
}

// Возвращает мапу с адресом эндпоинта и ролью, которая имеет доступ к нему
func (s *serv) accessibleRoles(ctx context.Context) (map[string]int64, error) {
	if accessibleRoles == nil {
		accessibleRoles = make(map[string]int64)

		// Лезем в базу за данными о доступных ролях для каждого эндпоинта
		// Можно кэшировать данные, чтобы не лезть в базу каждый раз

		// Например, для эндпоинта /note_v1.NoteV1/Get доступна только роль admin
		accessibleRoles[model.ExamplePath] = 1
	}

	return accessibleRoles, nil
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
