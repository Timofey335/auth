package converter

import (
	"time"

	modelCache "github.com/Timofey335/auth/internal/cache/user/model"
	"github.com/Timofey335/auth/internal/model"
)

// ToUserCacheFromUserModel - конвертер из модели UserModel в модель UserCacheModel
func ToUserCacheFromUserModel(user *model.UserModel) *modelCache.UserCacheModel {
	var updatedAt int64
	if !user.UpdatedAt.IsZero() {
		updatedAt = user.UpdatedAt.UnixNano()
	}

	return &modelCache.UserCacheModel{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.UnixNano(),
		UpdatedAt: &updatedAt,
	}
}

// ToUserModelFromUserCache - конвертер из модели UserCacheModel в модель UserModel
func ToUserModelFromUserCache(user *modelCache.UserCacheModel) *model.UserModel {
	createdAt := time.UnixMicro(user.CreatedAt)
	updatedAt := time.UnixMicro(*user.UpdatedAt)

	return &model.UserModel{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
