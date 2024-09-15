package cache

import (
	"context"

	"github.com/Timofey335/auth/internal/model"
)

// UserCache - интерфейс с методами пакета cache
type UserCache interface {
	CreateUser(ctx context.Context, user *model.UserModel) error
	GetUser(ctx context.Context, id int64) (*model.UserModel, error)
	DeleteUser(ctx context.Context, id int64) error
}
