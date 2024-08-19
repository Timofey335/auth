package repository

import (
	"context"

	desc "github.com/Timofey335/auth/pkg/auth_v1"
)

type UsersRepository interface {
	CreateUser(ctx context.Context, user *desc.CreateUserRequest) (int64, error)
	GetUser(ctx context.Context, userId int64) (*desc.GetUserResponse, error)
}
