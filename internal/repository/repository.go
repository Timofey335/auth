package repository

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Timofey335/auth/internal/model"
)

// UserRepository - интерфейс repo слоя
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.UserModel) (int64, error)
	GetUser(ctx context.Context, userId int64) (*model.UserModel, error)
	UpdateUser(ctx context.Context, user *model.UserUpdateModel) (*emptypb.Empty, error)
	DeleteUser(ctx context.Context, userId int64) (*emptypb.Empty, error)
}
