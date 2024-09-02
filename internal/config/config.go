package config

import (
	"github.com/joho/godotenv"
)

// GRPCConfig - интерфейс с методом Address
type GRPCConfig interface {
	Address() string
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
