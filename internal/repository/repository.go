package repository

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Timofey335/auth/internal/model"
	repoModel "github.com/Timofey335/auth/internal/repository/user/model"
)

// UserRepository - интерфейс repo слоя для User
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.UserModel) (int64, error)
	GetUser(ctx context.Context, userId int64) (*model.UserModel, error)
	UpdateUser(ctx context.Context, user *model.UserUpdateModel) (*emptypb.Empty, error)
	DeleteUser(ctx context.Context, userId int64) (*emptypb.Empty, error)
	GetUserData(ctx context.Context, userName string) (*repoModel.UserRepoModel, error)
}

// AccessRepository - интерфейс repo слоя для Access
type AccessRepository interface {
	GetRole(ctx context.Context, endpoint string) (int64, error)
}
