package env

import (
	"errors"
	"os"

	"github.com/Timofey335/auth/internal/config"
)

var _ config.PGConfig = (*pgConfig)(nil)

const (
	dsnEnvName = "PG_DSN"
)

type pgConfig struct {
	dsn string
}

// NewPGConfig - считывает значение DSN из env файла
func NewPGConfig() (*pgConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}

	return &pgConfig{
		dsn: dsn,
	}, nil
}

// DSN - возвращает из объекта pgConfig DSN для Postgres
func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
