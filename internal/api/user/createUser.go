package user

import (
	"context"

	"github.com/Timofey335/auth/internal/converter"
	desc "github.com/Timofey335/auth/pkg/auth_v1"
)

// CreateUser - создает нового пользователя
func (i *Implementation) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	id, err := i.userService.CreateUser(ctx, converter.ToUserFromDesc(req))
	if err != nil {
		return nil, err
	}

	return &desc.CreateUserResponse{
		Id: id,
	}, nil
}
