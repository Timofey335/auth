package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Timofey335/auth/internal/model"
	desc "github.com/Timofey335/auth/pkg/auth_v1"
)

func ToUserFromService(user *model.User) *desc.GetUserResponse {
	var updated_at *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updated_at = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.GetUserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      desc.Role(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updated_at,
	}
}

func ToUserFromDesc(user *desc.CreateUserRequest) *model.User {
	return &model.User{
		Name:            user.Name,
		Email:           user.Email,
		Password:        user.Password,
		PasswordConfirm: user.PasswordConfirm,
		Role:            int64(user.Role),
	}
}

func ToUserFromDescUpd(user *desc.UpdateUserRequest) *model.User {
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

	return &model.User{
		ID:              user.Id,
		Name:            name,
		Password:        password,
		PasswordConfirm: passwordConfirm,
		Role:            role,
	}
}
