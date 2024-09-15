package user

import (
	"github.com/Timofey335/platform_common/pkg/db"

	"github.com/Timofey335/auth/internal/cache"
	"github.com/Timofey335/auth/internal/repository"
	"github.com/Timofey335/auth/internal/service"
	def "github.com/Timofey335/auth/internal/service"
)

var _ def.UserService = (*serv)(nil)

type serv struct {
	userRepository repository.UserRepository
	cache          cache.UserCache
	txManager      db.TxManager
}

// NewService - создает новый экземпляр serv
func NewService(userRepository repository.UserRepository, userCache cache.UserCache, txManager db.TxManager) *serv {
	return &serv{
		userRepository: userRepository,
		cache:          userCache,
		txManager:      txManager,
	}
}

// NewMockService - mock для экземпляра serv
func NewMockService(deps ...interface{}) service.UserService {
	srv := serv{}

	for _, v := range deps {
		switch s := v.(type) {
		case repository.UserRepository:
			srv.userRepository = s
		}
	}

	return &srv
}
