package user

import (
	"github.com/Timofey335/auth/internal/repository"
	def "github.com/Timofey335/auth/internal/service"
)

var _ def.UserService = (*serv)(nil)

type serv struct {
	userRepository repository.UserRepository
}

func NewService(userRepository repository.UserRepository) *serv {
	return &serv{
		userRepository: userRepository,
	}
}
