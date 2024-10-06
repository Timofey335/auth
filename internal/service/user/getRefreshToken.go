package user

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Timofey335/auth/internal/model"
	"github.com/Timofey335/auth/internal/utils"
)

func (s *serv) GetRefreshToken(ctx context.Context, token string) (string, error) {
	refreshTokenSecretKey := s.authConfig.RefreshTokenSecretKey()
	refreshTokenExpiration := s.authConfig.RefreshTokenExpiration()

	claims, err := utils.VerifyToken(token, []byte(refreshTokenSecretKey))
	if err != nil {
		return "", status.Errorf(codes.Aborted, "invalid refresh token")
	}

	user, err := s.userRepository.GetUserData(ctx, claims.Username)
	if err != nil {
		return "", err
	}

	refreshToken, err := utils.GenerateToken(model.UserData{
		Username: user.Name,
		Role:     user.Role,
	},
		[]byte(refreshTokenSecretKey),
		time.Duration(refreshTokenExpiration),
	)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return refreshToken, nil
}
