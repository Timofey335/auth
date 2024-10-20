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

// GetRefreshToken - возвращает refresh токен
func (s *serv) GetRefreshToken(ctx context.Context, token string) (string, error) {
	refreshTokenSecretKey := s.authConfig.RefreshTokenSecretKey()
	refreshTokenExpiration := time.Duration(s.authConfig.RefreshTokenExpiration() * int64(time.Minute))

	claims, err := utils.VerifyToken(token, []byte(refreshTokenSecretKey))
	if err != nil {
		return "", status.Errorf(codes.Aborted, "invalid refresh token")
	}

	user, err := s.userRepository.GetUserData(ctx, claims.Email)
	if err != nil {
		return "", err
	}

	refreshToken, err := utils.GenerateToken(model.UserLoginModel{
		Email: claims.Email,
		Role:  user.Role,
	},
		[]byte(refreshTokenSecretKey),
		refreshTokenExpiration,
	)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return refreshToken, nil
}
