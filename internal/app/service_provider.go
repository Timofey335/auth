package app

import (
	"context"
	"log"

	"github.com/Timofey335/auth/internal/api/user"
	"github.com/Timofey335/auth/internal/client/db"
	"github.com/Timofey335/auth/internal/client/db/pg"
	"github.com/Timofey335/auth/internal/client/db/transaction"
	"github.com/Timofey335/auth/internal/closer"
	"github.com/Timofey335/auth/internal/config"
	"github.com/Timofey335/auth/internal/config/env"
	"github.com/Timofey335/auth/internal/repository"
	userRepository "github.com/Timofey335/auth/internal/repository/user"
	"github.com/Timofey335/auth/internal/service"
	userService "github.com/Timofey335/auth/internal/service/user"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient       db.Client
	txManager      db.TxManager
	userRepository repository.UserRepository

	userService service.UserService

	servImplementation *user.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// PGConfig - инициализирует конфигурацию БД из env файла
func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

// GRPCConfig - инициализирует конфигурацию GRPC сервера
func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

// DBClient - инициализирует подключение к БД
func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

// TxManager - инициализирует Transaction Manager
func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

// UserRepository - инициализация repo слоя
func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

// UserService - инициализация сервисного слоя
func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.userService
}

// ServImplementation - инициализация api слоя
func (s *serviceProvider) ServImplementation(ctx context.Context) *user.Implementation {
	if s.servImplementation == nil {
		s.servImplementation = user.NewImplementation(s.UserService(ctx))
	}

	return s.servImplementation
}
