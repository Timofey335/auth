package env

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
)

const (
	refreshTokenSecretKeyEnv  = "REFRESH_TOKEN_SECRET_KEY"
	refreshTokenExpirationEnv = "REFRESH_TOKEN_EXPIRATION"
)

type authSecretsConfig struct {
	refreshTokenSecretKey  string
	refreshTokenExpiration int64
}

// NewAuthConfig - считывает из env данные для генерации токенов
func NewAuthConfig() (*authSecretsConfig, error) {
	refreshTokenSecretKey := os.Getenv(refreshTokenSecretKeyEnv)
	if len(refreshTokenSecretKey) == 0 {
		return nil, errors.New("refresh token secret key not found")
	}

	refreshTokenExpirationStr := os.Getenv(refreshTokenExpirationEnv)
	if len(refreshTokenExpirationStr) == 0 {
		return nil, errors.New("refresh token expiration not found")
	}

	refreshTokenExpiration, err := strconv.ParseInt(refreshTokenExpirationStr, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse refresh token expiration")
	}

	return &authSecretsConfig{
		refreshTokenSecretKey:  refreshTokenSecretKey,
		refreshTokenExpiration: refreshTokenExpiration,
	}, nil
}

// RefreshTokenSecretKey - возвращет ключ для генерации токена
func (cfg *authSecretsConfig) RefreshTokenSecretKey() string {
	return cfg.refreshTokenSecretKey
}

// RefreshTokenExpiration - возвращет значение длительности срока годности токена
func (cfg *authSecretsConfig) RefreshTokenExpiration() int64 {
	return cfg.refreshTokenExpiration
}
