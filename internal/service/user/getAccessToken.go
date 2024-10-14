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

const (
	accessTokenSecretKey = "VqvguGiffXILza1f44TWXowDT4zwf03dtXmqWW4SYyE="
)

func (s *serv) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	refreshTokenSecretKey := s.authConfig.RefreshTokenSecretKey()
	// refreshTokenExpiration := s.authConfig.RefreshTokenExpiration()

	claims, err := utils.VerifyToken(refreshToken, []byte(refreshTokenSecretKey))
	if err != nil {
		return "", status.Errorf(codes.Aborted, "invalid refresh token")
	}

	user, err := s.userRepository.GetUserData(ctx, claims.Email)
	if err != nil {
		return "", err
	}

	accessToken, err := utils.GenerateToken(model.UserLoginModel{
		Email: user.Email,
		Role:  user.Role,
	},
		[]byte(accessTokenSecretKey),
		60*time.Minute,
	)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return accessToken, nil
}
