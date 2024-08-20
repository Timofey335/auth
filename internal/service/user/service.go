package user

import (
	"github.com/Timofey335/auth/internal/repository"
	def "github.com/Timofey335/auth/internal/service"
)

var _ def.UserService = (*serv)(nil)

type serv struct {
	userRepository repository.UsersRepository
}

func NewService(userRepository repository.UsersRepository) *serv {
	return &serv{
		userRepository: userRepository,
	}
}
