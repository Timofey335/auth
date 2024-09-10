package converter

import (
	"fmt"

	modelCache "github.com/Timofey335/auth/internal/cache/user/model"
	"github.com/Timofey335/auth/internal/model"
)

func ToUserCacheFromUserModel(user *model.UserModel) *modelCache.UserCacheModel {
	var updatedAt *int64
	if !user.UpdatedAt.IsZero() {
		*updatedAt = user.UpdatedAt.Unix()
	}
	fmt.Println(user.CreatedAt)
	return &modelCache.UserCacheModel{
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		Password:        user.Password,
		PasswordConfirm: user.PasswordConfirm,
		Role:            user.Role,
		CreatedAt:       user.CreatedAt.UnixNano(),
		UpdatedAt:       new(int64),
	}
}
