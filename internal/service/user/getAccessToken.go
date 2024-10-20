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

// GetAccessToken - возврещает access токен
func (s *serv) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	accessTokenSecretKey := s.authConfig.AccessTokenSecretKey()
	accessTokenExpiration := time.Duration(s.authConfig.AccessTokenExpiration() * int64(time.Minute))
	refreshTokenSecretKey := s.authConfig.RefreshTokenSecretKey()

	claims, err := utils.VerifyToken(refreshToken, []byte(refreshTokenSecretKey))
	if err != nil {
		return "", status.Errorf(codes.Aborted, "invalid refresh token")
	}

	user, err := s.userRepository.GetUserData(ctx, claims.Email)
	if err != nil {
		return "", err
	}

	accessToken, err := utils.GenerateToken(model.UserLoginModel{
		Email: claims.Email,
		Role:  user.Role,
	},
		[]byte(accessTokenSecretKey),
		accessTokenExpiration,
	)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return accessToken, nil
}
