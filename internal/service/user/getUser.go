package user

import (
	"context"
	"log"

	"github.com/fatih/color"

	"github.com/Timofey335/auth/internal/model"
)

// GetUser - получает данные о пользователе
func (s *serv) GetUser(ctx context.Context, id int64) (*model.UserModel, error) {
	var user *model.UserModel

	user, err := s.cache.GetUser(ctx, id)
	if err != nil {
		log.Println(color.HiMagentaString("getting the user with id-%d error: %v", id, err))

		user, err = s.userRepository.GetUser(ctx, id)
		if err != nil {
			return nil, err
		}

		err = s.cache.CreateUser(ctx, &model.UserModel{
			ID:        id,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}
