package user

import (
	"context"

	"github.com/Timofey335/auth/internal/model"
)

func (s *serv) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	id, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}
