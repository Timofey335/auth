package repository

import (
	"context"

	"github.com/Timofey335/auth/internal/model"
)

type UsersRepository interface {
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	GetUser(ctx context.Context, userId int64) (*model.User, error)
}

// type UsersRepository interface {
// 	CreateUser(ctx context.Context, user *model.User) (int64, error)
// 	GetUser(ctx context.Context, userId int64) (*model.User, error)
// }
