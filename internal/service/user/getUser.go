package user

import (
	"context"

	"github.com/Timofey335/auth/internal/model"
)

// GetUser - получает данные о пользователе
func (s *serv) GetUser(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.userRepository.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
