package access

import (
	"github.com/Timofey335/platform_common/pkg/db"

	"github.com/Timofey335/auth/internal/cache"
	"github.com/Timofey335/auth/internal/config"
	"github.com/Timofey335/auth/internal/repository"
	def "github.com/Timofey335/auth/internal/service"
)

var _ def.AccessService = (*serv)(nil)

type serv struct {
	accessRepository repository.AccessRepository
	userRepository   repository.UserRepository
	cache            cache.UserCache
	txManager        db.TxManager
	authConfig       config.AuthConfig
}

// NewService - создает новый экземпляр serv
func NewService(
	accessRepository repository.AccessRepository,
	userRepository repository.UserRepository,
	userCache cache.UserCache,
	txManager db.TxManager,
	authConfig config.AuthConfig,
) *serv {
	return &serv{
		accessRepository: accessRepository,
		userRepository:   userRepository,
		cache:            userCache,
		txManager:        txManager,
		authConfig:       authConfig,
	}
}
