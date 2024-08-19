package user

import (
	"github.com/Timofey335/auth/internal/repository"
	"github.com/Timofey335/auth/internal/service"
)

type serv struct {
	userRepository repository.UsersRepository
}

func NewService(userRepository repository.UsersRepository) service.UserService {
	return &serv{
		userRepository: userRepository,
	}
}
