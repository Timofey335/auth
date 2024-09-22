package config

import (
	"time"

	"github.com/joho/godotenv"
)

// GRPCConfig - интерфейс с методом Address
type GRPCConfig interface {
	Address() string
}

// HTTPConfig - интерфейс с методом Address
type HTTPConfig interface {
	Address() string
}

// RedisConfig - интерфейс конфигурации redis
type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
	UserExpiration() int64
}

// PGConfig - интерфейс с методом DSN
type PGConfig interface {
	DSN() string
}

// Load - считывает переменные из env файла
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
