package env

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
)

const (
	accessTokenSecretKeyEnv   = "ACCESS_TOKEN_SECRET_KEY"
	accessTokenExpirationEnv  = "ACCESS_TOKEN_EXPIRATION"
	refreshTokenSecretKeyEnv  = "REFRESH_TOKEN_SECRET_KEY"
	refreshTokenExpirationEnv = "REFRESH_TOKEN_EXPIRATION"
)

type authSecretsConfig struct {
	accesTokenSecretKey    string
	accessTokenExpiration  int64
	refreshTokenSecretKey  string
	refreshTokenExpiration int64
}

// NewAuthConfig - считывает из env данные для генерации токенов
func NewAuthConfig() (*authSecretsConfig, error) {
	accessTokenSecretKey := os.Getenv(accessTokenSecretKeyEnv)
	if len(accessTokenSecretKey) == 0 {
		return nil, errors.New("access token secret key not found")
	}

	accessTokenExpirationStr := os.Getenv(accessTokenExpirationEnv)
	if len(accessTokenExpirationStr) == 0 {
		return nil, errors.New("access token expiration not found")
	}

	accessTokenExpiration, err := strconv.ParseInt(accessTokenExpirationStr, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse access token expiration")
	}

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
		accesTokenSecretKey:    accessTokenSecretKey,
		accessTokenExpiration:  accessTokenExpiration,
		refreshTokenSecretKey:  refreshTokenSecretKey,
		refreshTokenExpiration: refreshTokenExpiration,
	}, nil
}

// AccessTokenSecretKey - возвращет ключ для генерации access токена
func (cfg *authSecretsConfig) AccessTokenSecretKey() string {
	return cfg.accesTokenSecretKey
}

// AccessTokenExpiration - возвращет значение длительности срока годности access токена
func (cfg *authSecretsConfig) AccessTokenExpiration() int64 {
	return cfg.accessTokenExpiration
}

// RefreshTokenSecretKey - возвращет ключ для генерации referesh токена
func (cfg *authSecretsConfig) RefreshTokenSecretKey() string {
	return cfg.refreshTokenSecretKey
}

// RefreshTokenExpiration - возвращет значение длительности срока годности refersh токена
func (cfg *authSecretsConfig) RefreshTokenExpiration() int64 {
	return cfg.refreshTokenExpiration
}
