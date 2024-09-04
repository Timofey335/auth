package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Timofey335/auth/internal/model"
	desc "github.com/Timofey335/auth/pkg/auth_v1"
)

// ToUserFromService - конвертирует данные из сервисного слоя для desc (GRPC)
func ToUserFromService(user *model.UserModel) *desc.GetUserResponse {
	var updatedAt *timestamppb.Timestamp
	if !user.UpdatedAt.IsZero() {
		updatedAt = timestamppb.New(user.UpdatedAt)
	}

	return &desc.GetUserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      desc.Role(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

// ToUserFromDesc - конвертирует данные из desc (GRPC) для сервисного слоя
func ToUserFromDesc(user *desc.CreateUserRequest) *model.UserModel {
	return &model.UserModel{
		Name:            user.Name,
		Email:           user.Email,
		Password:        user.Password,
		PasswordConfirm: user.PasswordConfirm,
		Role:            int64(user.Role),
	}
}

// ToUserFromDescUpd - конвертирует данные из desc (GRPC) для сервисного слоя
func ToUserFromDescUpd(user *desc.UpdateUserRequest) *model.UserUpdateModel {
	var name, password, passwordConfirm string
	var role int64

	if user.Name != nil {
		name = user.Name.Value
	}

	if user.Password != nil {
		password = user.Password.Value
	}

	if user.PasswordConfirm != nil {
		passwordConfirm = user.PasswordConfirm.Value
	}

	if user.Role != nil {
		switch *user.Role {
		case 1:
			role = 1
		case 2:
			role = 2
		default:
			role = 0
		}
	}

	return &model.UserUpdateModel{
		ID:              user.Id,
		Name:            &name,
		Password:        &password,
		PasswordConfirm: &passwordConfirm,
		Role:            &role,
	}
}
