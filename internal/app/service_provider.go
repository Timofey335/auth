package app

import (
	"context"
	"log"

	"github.com/Timofey335/auth/internal/api/user"
	"github.com/Timofey335/auth/internal/client/db"
	"github.com/Timofey335/auth/internal/client/db/pg"
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

	// pgPool         *pgxpool.Pool
	dbClient       db.Client
	userRepository repository.UserRepository

	userService service.UserService

	servImplementation *user.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

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

// func (s *serviceProvider) PgPool(ctx context.Context) *pgxpool.Pool {
// 	if s.pgPool == nil {
// 		pool, err := pgxpool.Connect(ctx, s.PGConfig().DSN())
// 		if err != nil {
// 			log.Fatalf("failed to connect to database: %v", err)
// 		}

// 		if err = pool.Ping(ctx); err != nil {
// 			log.Fatalf("ping error: %s", err.Error())
// 		}

// 		closer.Add(func() error {
// 			pool.Close()
// 			return nil
// 		})

// 		s.pgPool = pool
// 	}

// 	return s.pgPool
// }

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.UserRepository(ctx))
	}

	return s.userService
}

func (s *serviceProvider) ServImplementation(ctx context.Context) *user.Implementation {
	if s.servImplementation == nil {
		s.servImplementation = user.NewImplementation(s.UserService(ctx))
	}

	return s.servImplementation
}
