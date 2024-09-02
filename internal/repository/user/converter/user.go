package converter

import (
	"github.com/Timofey335/auth/internal/model"
	modelRepo "github.com/Timofey335/auth/internal/repository/user/model"
)

// ToUserFromRepo - конвертирует данные из repo слоя в сервисный слой
func ToUserFromRepo(user *modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
