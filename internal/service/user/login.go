package user

import (
	"context"
	"errors"
	"time"

	"github.com/Timofey335/auth/internal/model"
	"github.com/Timofey335/auth/internal/utils"
)

// Login - метод для аутентификации и получения токена
func (s *serv) Login(ctx context.Context, userLoginData *model.UserLoginModel) (string, error) {
	user, err := s.userRepository.GetUserData(ctx, userLoginData.Email)
	if err != nil {
		return "", err
	}

	if !utils.VerifyPassword(user.Password, userLoginData.Password) {
		return "", errors.New("user not found or password incorrect")
	}

	refreshTokenExpiration := time.Duration(s.authConfig.RefreshTokenExpiration() * int64(time.Minute))

	refreshToken, err := utils.GenerateToken(model.UserData{
		Username: userLoginData.Email,
		Role:     user.Role,
	},
		[]byte(s.authConfig.RefreshTokenSecretKey()),
		refreshTokenExpiration,
	)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return refreshToken, nil
}
