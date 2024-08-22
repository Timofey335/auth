package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Timofey335/auth/internal/model"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	GetUser(ctx context.Context, userId int64) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) (*emptypb.Empty, error)
	DeleteUser(ctx context.Context, userId int64) (*emptypb.Empty, error)
}
