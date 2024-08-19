package service

import (
	"context"

	"github.com/Timofey335/auth/internal/model"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	GetUser(ctx context.Context, id int64) (*model.User, error)
}
